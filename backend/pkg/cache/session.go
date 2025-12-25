package cache

import (
	"fmt"
	"time"
)

const (
	// SessionTTL is how long sessions are cached (matches JWT expiry)
	SessionTTL = 24 * time.Hour

	// SessionPrefix is the Redis key prefix for sessions
	SessionPrefix = "session:"

	// UserSessionsPrefix is the Redis key prefix for user's sessions set
	UserSessionsPrefix = "user_sessions:"
)

// Session represents a cached user session
type Session struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	CreatedAt    int64  `json:"createdAt"`
	ExpiresAt    int64  `json:"expiresAt"`
	IP           string `json:"ip,omitempty"`
	UserAgent    string `json:"userAgent,omitempty"`
}

// SetSession caches a user session
func SetSession(sessionID string, session *Session, ttl time.Duration) error {
	key := SessionPrefix + sessionID
	if ttl == 0 {
		ttl = SessionTTL
	}

	// Store session
	if err := SetJSON(key, session, ttl); err != nil {
		return err
	}

	// Add to user's sessions set
	userSessionsKey := UserSessionsPrefix + session.UserID
	_ = SAdd(userSessionsKey, sessionID)
	_ = Expire(userSessionsKey, ttl)

	return nil
}

// GetSession retrieves a cached session
func GetSession(sessionID string) (*Session, error) {
	key := SessionPrefix + sessionID
	session := &Session{}
	err := GetJSON(key, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// DeleteSession removes a session
func DeleteSession(sessionID string) error {
	// Get session first to find user ID
	session, err := GetSession(sessionID)
	if err == nil && session != nil {
		// Remove from user's sessions set
		userSessionsKey := UserSessionsPrefix + session.UserID
		_ = SRem(userSessionsKey, sessionID)
	}

	key := SessionPrefix + sessionID
	return Delete(key)
}

// SessionExists checks if a session exists
func SessionExists(sessionID string) (bool, error) {
	key := SessionPrefix + sessionID
	return Exists(key)
}

// GetUserSessions returns all session IDs for a user
func GetUserSessions(userID string) ([]string, error) {
	key := UserSessionsPrefix + userID
	return SMembers(key)
}

// DeleteAllUserSessions removes all sessions for a user (for logout all devices)
func DeleteAllUserSessions(userID string) error {
	sessions, err := GetUserSessions(userID)
	if err != nil {
		return err
	}

	for _, sessionID := range sessions {
		_ = Delete(SessionPrefix + sessionID)
	}

	// Remove the user sessions set
	return Delete(UserSessionsPrefix + userID)
}

// CountUserSessions returns the number of active sessions for a user
func CountUserSessions(userID string) (int64, error) {
	key := UserSessionsPrefix + userID
	members, err := SMembers(key)
	if err != nil {
		return 0, err
	}
	return int64(len(members)), nil
}

// RefreshSession extends the session TTL
func RefreshSession(sessionID string, ttl time.Duration) error {
	key := SessionPrefix + sessionID
	if ttl == 0 {
		ttl = SessionTTL
	}
	return Expire(key, ttl)
}

// ValidateSession checks if session exists and is not expired
func ValidateSession(sessionID string) (bool, error) {
	session, err := GetSession(sessionID)
	if err != nil {
		return false, err
	}
	if session == nil {
		return false, nil
	}

	// Check if expired
	if session.ExpiresAt > 0 && session.ExpiresAt < time.Now().UnixMilli() {
		_ = DeleteSession(sessionID)
		return false, fmt.Errorf("session expired")
	}

	return true, nil
}

// SessionStats returns statistics about cached sessions
func SessionStats() map[string]interface{} {
	sessionKeys, _ := Keys(SessionPrefix + "*")
	userSessionKeys, _ := Keys(UserSessionsPrefix + "*")

	return map[string]interface{}{
		"activeSessions": len(sessionKeys),
		"activeUsers":    len(userSessionKeys),
		"connected":      IsConnected(),
	}
}
