package dto

type (
	PositionResponse struct {
		ID               string  `json:"id"`
		PortfolioID      string  `json:"portfolioId"`
		InstrumentID     string  `json:"instrumentId"`
		Symbol           string  `json:"symbol"`
		Quantity         float64 `json:"quantity"`
		AvgCost          float64 `json:"avgCost"`
		TotalCost        float64 `json:"totalCost"`
		CurrentPrice     float64 `json:"currentPrice,omitempty"`
		MarketValue      float64 `json:"marketValue,omitempty"`
		UnrealizedPnL    float64 `json:"unrealizedPnL,omitempty"`
		UnrealizedPnLPct float64 `json:"unrealizedPnLPct,omitempty"`
	}

	PositionLotResponse struct {
		ID           string  `json:"id"`
		Quantity     float64 `json:"quantity"`
		RemainingQty float64 `json:"remainingQty"`
		CostPerUnit  float64 `json:"costPerUnit"`
		PurchasedAt  string  `json:"purchasedAt"`
	}

	PositionDetailResponse struct {
		Position PositionResponse      `json:"position"`
		Lots     []PositionLotResponse `json:"lots"`
	}
)
