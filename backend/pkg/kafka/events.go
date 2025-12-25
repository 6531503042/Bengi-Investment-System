package kafka

import (
	"encoding/json"
	"log"
	"time"
)

// Trade Event types
const (
	EventTradeExecuted    = "TRADE_EXECUTED"
	EventPortfolioUpdated = "PORTFOLIO_UPDATED"
	EventAccountUpdated   = "ACCOUNT_UPDATED"
)

// TradeEvent represents a trade execution event
type TradeEvent struct {
	EventType   string  `json:"eventType"`
	TradeID     string  `json:"tradeId"`
	OrderID     string  `json:"orderId"`
	UserID      string  `json:"userId"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Commission  float64 `json:"commission"`
	TotalAmount float64 `json:"totalAmount"`
	Timestamp   int64   `json:"timestamp"`
}

// PortfolioEvent represents a portfolio update event
type PortfolioEvent struct {
	EventType    string  `json:"eventType"`
	PortfolioID  string  `json:"portfolioId"`
	UserID       string  `json:"userId"`
	Symbol       string  `json:"symbol"`
	Action       string  `json:"action"` // ADD, UPDATE, REMOVE
	Quantity     float64 `json:"quantity"`
	AverageCost  float64 `json:"averageCost"`
	CurrentValue float64 `json:"currentValue"`
	Timestamp    int64   `json:"timestamp"`
}

// AccountEvent represents an account balance update
type AccountEvent struct {
	EventType       string  `json:"eventType"`
	AccountID       string  `json:"accountId"`
	UserID          string  `json:"userId"`
	TransactionType string  `json:"transactionType"`
	Amount          float64 `json:"amount"`
	NewBalance      float64 `json:"newBalance"`
	Currency        string  `json:"currency"`
	Timestamp       int64   `json:"timestamp"`
}

// PublishTradeExecuted publishes a trade execution event
func PublishTradeExecuted(trade *TradeEvent) error {
	trade.EventType = EventTradeExecuted
	trade.Timestamp = time.Now().UnixMilli()

	if !IsConnected() {
		log.Printf("[Kafka] Not connected, skipping trade event")
		return nil
	}

	data, err := json.Marshal(trade)
	if err != nil {
		return err
	}

	log.Printf("[Kafka] Publishing trade event: %s for %s", trade.TradeID, trade.Symbol)
	return Publish(TopicTrades, trade.TradeID, data)
}

// PublishPortfolioUpdated publishes a portfolio update event
func PublishPortfolioUpdated(event *PortfolioEvent) error {
	event.EventType = EventPortfolioUpdated
	event.Timestamp = time.Now().UnixMilli()

	if !IsConnected() {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return Publish(TopicEvents, event.PortfolioID, data)
}

// PublishAccountUpdated publishes an account update event
func PublishAccountUpdated(event *AccountEvent) error {
	event.EventType = EventAccountUpdated
	event.Timestamp = time.Now().UnixMilli()

	if !IsConnected() {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return Publish(TopicEvents, event.AccountID, data)
}

// StartTradeConsumer starts consuming trade events
func StartTradeConsumer(handler func(event *TradeEvent) error) {
	go func() {
		ConsumeTrades(func(key, value []byte) error {
			var event TradeEvent
			if err := json.Unmarshal(value, &event); err != nil {
				log.Printf("[Kafka] Failed to unmarshal trade event: %v", err)
				return err
			}
			return handler(&event)
		})
	}()
}
