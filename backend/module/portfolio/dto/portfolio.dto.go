package dto

type (
	CreatePortfolioRequest struct {
		AccountID string `json:"accountId" validate:"required"`
		Name      string `json:"name" validate:"required,min=1,max=50"`
	}

	UpdatePortfolioRequest struct {
		Name      string `json:"name" validate:"omitempty,min=1,max=50"`
		IsDefault bool   `json:"isDefault"`
	}

	PortfolioResponse struct {
		ID        string  `json:"id"`
		UserID    string  `json:"userId"`
		AccountID string  `json:"accountId"`
		Name      string  `json:"name"`
		IsDefault bool    `json:"isDefault"`
		Value     float64 `json:"value,omitempty"`
	}

	PortfolioSummary struct {
		Portfolio   PortfolioResponse  `json:"portfolio"`
		Positions   []PositionResponse `json:"positions"`
		TotalValue  float64            `json:"totalValue"`
		TotalCost   float64            `json:"totalCost"`
		TotalPnL    float64            `json:"totalPnL"`
		TotalPnLPct float64            `json:"totalPnLPct"`
	}
)
