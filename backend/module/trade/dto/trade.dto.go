package dto

type (
	ExecuteTradeRequest struct {
		OrderID  string  `json:"orderId" validate:"required"`
		Price    float64 `json:"price" validate:"required, gt=0"`
		Quantity float64 `json:"quantity" validate:"required, gt=0"`
	}

	TradeResponse struct {
		ID           string  `json:"id"`
		OrderID      string  `json:"orderId"`
		UserID       string  `json:"userId"`
		AccountID    string  `json:"accountId"`
		PortfolioID  string  `json:"portfolioId"`
		InstrumentID string  `json:"instrumentId"`
		Symbol       string  `json:"symbol"`
		Side         string  `json:"side"`
		Quantity     float64 `json:"quantity"`
		Price        float64 `json:"price"`
		Total        float64 `json:"total"`
		Commission   float64 `json:"commission"`
		NetAmount    float64 `json:"netAmount"`
		ExecutedAt   string  `json:"executedAt"`
	}

	TradeListResponse struct {
		Trades []TradeResponse `json:"trades"`
		Total  int             `json:"total"`
		Page   int             `json:"page"`
		Limit  int             `json:"limit"`
	}

	TradeFilter struct {
		Symbol string `query:"symbol"`
		Side   string `query:"side"`
		Page   int    `query:"page"`
		Limit  int    `query:"limit"`
	}

	TradeSummary struct {
		TotalTrades     int     `json:"totalTrades"`
		TotalValue      int     `json:"totalValue"`
		TotalBuyValue   float64 `json:"totalBuyValue"`
		TotalSellValue  float64 `json:"totalSellValue"`
		TotalCommission float64 `json:"totalCommission"`
		NetProfit       float64 `json:"netProfit"`
	}
)
