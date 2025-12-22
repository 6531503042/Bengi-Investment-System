package service

import (
	"context"
	"errors"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MaxWatchlistsPerUser = 10
const MaxSymbolsPerWatchlist = 50

var (
	ErrWatchlistNotFound    = errors.New("watchlist not found")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrMaxWatchlistsReached = errors.New("maximum watchlists reached")
	ErrMaxSymbolsReached    = errors.New("maximum symbols reached")
	ErrSymbolAlreadyExists  = errors.New("symbol already in watchlist")
	ErrSymbolNotFound       = errors.New("symbol not in watchlist")
)

type WatchlistService struct {
	repo *repository.WatchlistRepository
}

func NewWatchlistService(repo *repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{repo: repo}
}

// CreateWatchlist creates a new watchlist
func (s *WatchlistService) CreateWatchlist(ctx context.Context, userID string, req *dto.CreateWatchlistRequest) (*dto.WatchlistResponse, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Check max watchlists limit
	count, err := s.repo.CountByUserID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}
	if count >= MaxWatchlistsPerUser {
		return nil, ErrMaxWatchlistsReached
	}

	watchlist := &model.Watchlist{
		UserID:    userObjectID,
		Name:      req.Name,
		Symbols:   req.Symbols,
		IsDefault: count == 0, // First watchlist is default
	}

	if err := s.repo.Create(ctx, watchlist); err != nil {
		return nil, err
	}

	return s.toWatchlistResponse(watchlist), nil
}

// GetWatchlists returns all watchlists for a user
func (s *WatchlistService) GetWatchlists(ctx context.Context, userID string) ([]dto.WatchlistResponse, error) {
	watchlists, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.WatchlistResponse
	for _, w := range watchlists {
		responses = append(responses, *s.toWatchlistResponse(&w))
	}
	return responses, nil
}

// GetWatchlistByID returns a single watchlist
func (s *WatchlistService) GetWatchlistByID(ctx context.Context, watchlistID, userID string) (*dto.WatchlistResponse, error) {
	watchlist, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, ErrWatchlistNotFound
	}

	if watchlist.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	return s.toWatchlistResponse(watchlist), nil
}

// UpdateWatchlist updates a watchlist
func (s *WatchlistService) UpdateWatchlist(ctx context.Context, watchlistID, userID string, req *dto.UpdateWatchlistRequest) (*dto.WatchlistResponse, error) {
	watchlist, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, ErrWatchlistNotFound
	}

	if watchlist.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}

	if req.IsDefault {
		// Clear other defaults first
		if err := s.repo.ClearDefaultWatchlists(ctx, watchlist.UserID); err != nil {
			return nil, err
		}
		update["isDefault"] = true
	}

	if len(update) > 0 {
		if err := s.repo.Update(ctx, watchlist.ID, update); err != nil {
			return nil, err
		}
	}

	updated, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, err
	}

	return s.toWatchlistResponse(updated), nil
}

// DeleteWatchlist deletes a watchlist
func (s *WatchlistService) DeleteWatchlist(ctx context.Context, watchlistID, userID string) error {
	watchlist, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return ErrWatchlistNotFound
	}

	if watchlist.UserID.Hex() != userID {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, watchlist.ID)
}

// AddSymbol adds a symbol to watchlist
func (s *WatchlistService) AddSymbol(ctx context.Context, watchlistID, userID string, req *dto.AddSymbolRequest) (*dto.WatchlistResponse, error) {
	watchlist, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, ErrWatchlistNotFound
	}

	if watchlist.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	// Check max symbols limit
	if len(watchlist.Symbols) >= MaxSymbolsPerWatchlist {
		return nil, ErrMaxSymbolsReached
	}

	// Check if symbol already exists
	for _, sym := range watchlist.Symbols {
		if sym == req.Symbol {
			return nil, ErrSymbolAlreadyExists
		}
	}

	if err := s.repo.AddSymbol(ctx, watchlist.ID, req.Symbol); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, err
	}

	return s.toWatchlistResponse(updated), nil
}

// RemoveSymbol removes a symbol from watchlist
func (s *WatchlistService) RemoveSymbol(ctx context.Context, watchlistID, userID string, req *dto.RemoveSymbolRequest) (*dto.WatchlistResponse, error) {
	watchlist, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, ErrWatchlistNotFound
	}

	if watchlist.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	// Check if symbol exists
	found := false
	for _, sym := range watchlist.Symbols {
		if sym == req.Symbol {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrSymbolNotFound
	}

	if err := s.repo.RemoveSymbol(ctx, watchlist.ID, req.Symbol); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(ctx, watchlistID)
	if err != nil {
		return nil, err
	}

	return s.toWatchlistResponse(updated), nil
}

func (s *WatchlistService) toWatchlistResponse(w *model.Watchlist) *dto.WatchlistResponse {
	symbols := w.Symbols
	if symbols == nil {
		symbols = []string{}
	}

	return &dto.WatchlistResponse{
		ID:        w.ID.Hex(),
		UserID:    w.UserID.Hex(),
		Name:      w.Name,
		Symbols:   symbols,
		IsDefault: w.IsDefault,
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
		UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
	}
}
