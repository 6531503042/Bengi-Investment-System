package dto

// Order side constants - used across the trading platform
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// Order type constants - determines how the order is executed
type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET" // Execute immediately at current price
	OrderTypeLimit  OrderType = "LIMIT"  // Execute only at specified price or better
	OrderTypeStop   OrderType = "STOP"   // Trigger when price reaches stop price
)

// Time in force - how long the order remains active
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good 'til canceled
	TimeInForceDAY TimeInForce = "DAY" // Valid for current trading day only
	TimeInForceIOC TimeInForce = "IOC" // Immediate or cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or kill - all or nothing
)

// CreateOrderRequest contains the data needed to place a new order.
// All monetary values are in the account's base currency.
type CreateOrderRequest struct {
	AccountID   string  `json:"accountId" validate:"required"`
	PortfolioID string  `json:"portfolioId" validate:"required"`
	Symbol      string  `json:"symbol" validate:"required"`
	Side        string  `json:"side" validate:"required,oneof=BUY SELL"`
	Type        string  `json:"type" validate:"required,oneof=MARKET LIMIT STOP"`
	Quantity    float64 `json:"quantity" validate:"required,gt=0"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`     // Required for LIMIT orders
	StopPrice   float64 `json:"stopPrice" validate:"omitempty,gt=0"` // Required for STOP orders
	TimeInForce string  `json:"timeInForce" validate:"omitempty,oneof=GTC DAY IOC FOK"`
}

// OrderResponse represents a complete order with all its details.
// Used when returning order data to clients.
type OrderResponse struct {
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
	FilledQty    float64 `json:"filledQty"`              // Partially filled amount
	Price        float64 `json:"price,omitempty"`        // Limit price if applicable
	StopPrice    float64 `json:"stopPrice,omitempty"`    // Stop trigger price
	AvgFillPrice float64 `json:"avgFillPrice,omitempty"` // Weighted average of all fills
	Commission   float64 `json:"commission"`
	CreatedAt    string  `json:"createdAt"`
	FilledAt     *string `json:"filledAt,omitempty"`
	CancelledAt  *string `json:"cancelledAt,omitempty"`
}

// OrderListResponse wraps a paginated list of orders.
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

// OrderFilter allows filtering orders by various criteria.
type OrderFilter struct {
	Status string `query:"status"` // PENDING, FILLED, CANCELLED, etc.
	Side   string `query:"side"`   // BUY or SELL
	Symbol string `query:"symbol"` // Filter by trading symbol
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
}

// CancelOrderRequest contains the reason for cancelling an order.
type CancelOrderRequest struct {
	Reason string `json:"reason"`
}
