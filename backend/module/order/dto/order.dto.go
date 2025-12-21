package dto

type (
	CreateOrderRequest struct {
		AccountID   string  `json:"accountId" validate:"required"`
		PortfolioID string  `json:"portfolioId" validate:"required"`
		Symbol      string  `json:"symbol" validate:"required"`
		Side        string  `json:"side" validate:"required,oneof=BUY SELL"`
		Type        string  `json:"type" validate:"required,oneof=MARKET LIMIT STOP"`
		Quantity    float64 `json:"quantity" validate:"required,gt=0"`
		Price       float64 `json:"price" validate:"omitempty,gt=0"`
		StopPrice   float64 `json:"stopPrice" validate:"omitempty,gt=0"`
		TimeInForce string  `json:"timeInForce" validate:"omitempty,oneof=GTC DAY IOC FOK"`
	}

	OrderResponse struct {
		ID           string  `json:"id"`
		UserID       string  `json:"userId"`
		AccountID    string  `json:"accountId"`
		PortfolioID  string  `json:"portfolioId"`
		InstrumentID string  `json:"instrumentId"`
		Symbol       string  `json:"symbol"`
		Side         string  `json:"side"`
		Type         string  `json:"type"`
		Status       string  `json:"status"`
		TimeInForce  string  `json:"timeInForce"`
		Quantity     float64 `json:"quantity"`
		FilledQty    float64 `json:"filledQty"`
		Price        float64 `json:"price,omitempty"`
		StopPrice    float64 `json:"stopPrice,omitempty"`
		AvgFillPrice float64 `json:"avgFillPrice,omitempty"`
		Commission   float64 `json:"commission"`
		CreatedAt    string  `json:"createdAt"`
		FilledAt     *string `json:"filledAt,omitempty"`
		CancelledAt  *string `json:"cancelledAt,omitempty"`
	}

	OrderListResponse struct {
		Orders []OrderResponse `json:"orders"`
		Total  int             `json:"total"`
		Page   int             `json:"page"`
		Limit  int             `json:"limit"`
	}

	OrderFilter struct {
		Status string `query:"status"`
		Side   string `query:"side"`
		Symbol string `query:"symbol"`
		Page   int    `query:"page"`
		Limit  int    `query:"limit"`
	}

	CancelOrderRequest struct {
		Reason string `json:"reason"`
	}
)
