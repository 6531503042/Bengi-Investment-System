package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/account/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNotDemoAccount = errors.New("operation only allowed for demo accounts")
)

type DemoService struct {
	repository *repository.AccountRepository
}

func NewDemoService(repository *repository.AccountRepository) *DemoService {
	return &DemoService{
		repository: repository,
	}
}

// CreateDemoAccount creates a new demo account with virtual funds
func (s *DemoService) CreateDemoAccount(ctx context.Context, userID string, req *dto.CreateDemoAccountRequest) (*dto.CreateDemoAccountResponse, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Set defaults
	currency := "USD"
	if req.Currency != "" {
		currency = req.Currency
	}

	leverage := 10
	if req.Leverage > 0 && req.Leverage <= 100 {
		leverage = req.Leverage
	}

	initialBalance := 10000.0
	if req.InitialBalance > 0 {
		initialBalance = req.InitialBalance
	}

	account := &model.Account{
		UserID:         userObjectID,
		Currency:       currency,
		Balance:        initialBalance,
		Status:         model.AccountStatusActive,
		Type:           model.AccountTypeDemo,
		Leverage:       leverage,
		InitialBalance: initialBalance,
		TotalDeposits:  initialBalance,
		TotalPnL:       0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repository.Create(ctx, account); err != nil {
		return nil, err
	}

	return &dto.CreateDemoAccountResponse{
		AccountID:      account.ID.Hex(),
		Currency:       account.Currency,
		Balance:        account.Balance,
		Leverage:       account.Leverage,
		InitialBalance: account.InitialBalance,
		Message:        "Demo account created successfully with $" + formatFloat(account.Balance),
	}, nil
}

// DemoDeposit adds virtual funds to demo account (unlimited)
func (s *DemoService) DemoDeposit(ctx context.Context, accountID, userID string, req *dto.DemoDepositRequest) (*dto.DemoDepositResponse, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Check ownership
	if account.UserID.Hex() != userID {
		return nil, errors.New("account not found")
	}

	// Only demo accounts can use this
	if account.Type != model.AccountTypeDemo {
		return nil, ErrNotDemoAccount
	}

	// Add funds
	newBalance := account.Balance + req.Amount
	newTotalDeposits := account.TotalDeposits + req.Amount

	if err := s.repository.UpdateBalance(ctx, account.ID, newBalance); err != nil {
		return nil, err
	}

	// Update total deposits
	if err := s.repository.UpdateField(ctx, account.ID, "totalDeposits", newTotalDeposits); err != nil {
		return nil, err
	}

	return &dto.DemoDepositResponse{
		AccountID:     account.ID.Hex(),
		NewBalance:    newBalance,
		TotalDeposits: newTotalDeposits,
		Message:       "Deposited $" + formatFloat(req.Amount) + " to demo account",
	}, nil
}

// DemoReset resets demo account to initial state
func (s *DemoService) DemoReset(ctx context.Context, accountID, userID string, req *dto.DemoResetRequest) (*dto.DemoResetResponse, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Check ownership
	if account.UserID.Hex() != userID {
		return nil, errors.New("account not found")
	}

	// Only demo accounts can use this
	if account.Type != model.AccountTypeDemo {
		return nil, ErrNotDemoAccount
	}

	// Reset balance
	newBalance := account.InitialBalance
	if req.InitialBalance > 0 {
		newBalance = req.InitialBalance
	}

	if err := s.repository.UpdateBalance(ctx, account.ID, newBalance); err != nil {
		return nil, err
	}

	// Reset other fields
	_ = s.repository.UpdateField(ctx, account.ID, "totalDeposits", newBalance)
	_ = s.repository.UpdateField(ctx, account.ID, "totalPnL", 0.0)
	_ = s.repository.UpdateField(ctx, account.ID, "initialBalance", newBalance)

	return &dto.DemoResetResponse{
		AccountID:      account.ID.Hex(),
		NewBalance:     newBalance,
		InitialBalance: newBalance,
		Message:        "Demo account reset to $" + formatFloat(newBalance),
	}, nil
}

// GetDemoStats returns demo account statistics
func (s *DemoService) GetDemoStats(ctx context.Context, accountID, userID string) (*dto.DemoAccountStats, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Check ownership
	if account.UserID.Hex() != userID {
		return nil, errors.New("account not found")
	}

	// Only demo accounts
	if account.Type != model.AccountTypeDemo {
		return nil, ErrNotDemoAccount
	}

	pnlPercentage := 0.0
	if account.InitialBalance > 0 {
		pnlPercentage = ((account.Balance - account.InitialBalance) / account.InitialBalance) * 100
	}

	return &dto.DemoAccountStats{
		AccountID:      account.ID.Hex(),
		Currency:       account.Currency,
		Balance:        account.Balance,
		InitialBalance: account.InitialBalance,
		TotalDeposits:  account.TotalDeposits,
		TotalPnL:       account.TotalPnL,
		Leverage:       account.Leverage,
		PnLPercentage:  pnlPercentage,
		CreatedAt:      account.CreatedAt,
	}, nil
}

// GetOrCreateDemoAccount gets existing demo account or creates new one
func (s *DemoService) GetOrCreateDemoAccount(ctx context.Context, userID string) (*dto.DemoAccountStats, error) {
	accounts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find existing demo account
	for _, acc := range accounts {
		if acc.Type == model.AccountTypeDemo {
			return s.GetDemoStats(ctx, acc.ID.Hex(), userID)
		}
	}

	// Create new demo account if none exists
	resp, err := s.CreateDemoAccount(ctx, userID, &dto.CreateDemoAccountRequest{
		Currency:       "USD",
		Leverage:       10,
		InitialBalance: 10000,
	})
	if err != nil {
		return nil, err
	}

	return &dto.DemoAccountStats{
		AccountID:      resp.AccountID,
		Currency:       resp.Currency,
		Balance:        resp.Balance,
		InitialBalance: resp.InitialBalance,
		TotalDeposits:  resp.InitialBalance,
		TotalPnL:       0,
		Leverage:       resp.Leverage,
		PnLPercentage:  0,
		CreatedAt:      time.Now(),
	}, nil
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
