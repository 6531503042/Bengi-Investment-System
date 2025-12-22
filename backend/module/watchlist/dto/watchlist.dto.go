package dto

type (
	CreateWatchlistRequest struct {
		Name    string   `json:"name" validate:"required,min=1,max=50"`
		Symbols []string `json:"symbols"`
	}

	UpdateWatchlistRequest struct {
		Name      string `json:"name" validate:"omitempty,min=1,max=50"`
		IsDefault bool   `json:"isDefault"`
	}

	AddSymbolRequest struct {
		Symbol string `json:"symbol" validate:"required"`
	}

	RemoveSymbolRequest struct {
		Symbol string `json:"symbol" validate:"required"`
	}

	WatchlistResponse struct {
		ID        string   `json:"id"`
		UserID    string   `json:"userId"`
		Name      string   `json:"name"`
		Symbols   []string `json:"symbols"`
		IsDefault bool     `json:"isDefault"`
		CreatedAt string   `json:"createdAt"`
		UpdatedAt string   `json:"updatedAt"`
	}

	WatchlistWithQuotes struct {
		Watchlist WatchlistResponse `json:"watchlist"`
		Quotes    []QuoteInfo       `json:"quotes"`
	}

	QuoteInfo struct {
		Symbol        string  `json:"symbol"`
		Price         float64 `json:"price"`
		Change        float64 `json:"change"`
		ChangePercent float64 `json:"changePercent"`
	}
)
