package service

import (
	"context"
	"errors"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/auth/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/repository"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/cache"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists         = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrWrongPassword       = errors.New("current password is incorrect")
	ErrSessionNotFound     = errors.New("session not found or expired")
)

type AuthService struct {
	UserRepository *repository.UserRepository
}

func NewAuthService(userRepository *repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	exists, err := s.UserRepository.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		RoleID:   primitive.NewObjectID(),
	}

	if err := s.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		FullName: user.FullName,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.UserRepository.UpdateLastLogin(ctx, user.ID)

	accessToken, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.RoleID.Hex())
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Store session in Redis (if available)
	sessionID := uuid.New().String()
	if cache.IsConnected() {
		session := &cache.Session{
			UserID:       user.ID.Hex(),
			Email:        user.Email,
			Role:         user.RoleID.Hex(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			CreatedAt:    time.Now().UnixMilli(),
			ExpiresAt:    time.Now().Add(24 * time.Hour).UnixMilli(),
		}
		_ = cache.SetSession(sessionID, session, 24*time.Hour)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
		User: dto.UserProfile{
			ID:       user.ID.Hex(),
			Email:    user.Email,
			FullName: user.FullName,
			RoleID:   user.RoleID.Hex(),
		},
	}, nil
}

// Logout removes the session from Redis
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if cache.IsConnected() && sessionID != "" {
		return cache.DeleteSession(sessionID)
	}
	return nil
}

// LogoutAll removes all sessions for a user
func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	if cache.IsConnected() {
		return cache.DeleteAllUserSessions(userID)
	}
	return nil
}

// ValidateSession checks if a session is valid
func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (bool, error) {
	if !cache.IsConnected() {
		// If Redis not available, skip session validation
		return true, nil
	}
	return cache.ValidateSession(sessionID)
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*dto.UserProfile, error) {
	user, err := s.UserRepository.FindById(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfile{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   user.RoleID.Hex(),
	}, nil
}

// RefreshToken validates refresh token and generates new tokens
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Get user from database
	user, err := s.UserRepository.FindById(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.RoleID.Hex())
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Create new session
	sessionID := uuid.New().String()
	if cache.IsConnected() {
		session := &cache.Session{
			UserID:       user.ID.Hex(),
			Email:        user.Email,
			Role:         user.RoleID.Hex(),
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
			CreatedAt:    time.Now().UnixMilli(),
			ExpiresAt:    time.Now().Add(24 * time.Hour).UnixMilli(),
		}
		_ = cache.SetSession(sessionID, session, 24*time.Hour)
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		SessionID:    sessionID,
		User: dto.UserProfile{
			ID:       user.ID.Hex(),
			Email:    user.Email,
			FullName: user.FullName,
			RoleID:   user.RoleID.Hex(),
		},
	}, nil
}

// ChangePassword changes user password after validating current password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.UserRepository.FindById(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		return ErrWrongPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.UserRepository.UpdatePassword(ctx, user.ID, string(hashedPassword))
}
