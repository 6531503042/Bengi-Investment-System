package kafka

import (
	"encoding/json"
	"log"
	"time"
)

// OrderEvent types
const (
	EventOrderCreated   = "ORDER_CREATED"
	EventOrderCancelled = "ORDER_CANCELLED"
	EventOrderFilled    = "ORDER_FILLED"
	EventOrderPartial   = "ORDER_PARTIALLY_FILLED"
)

// OrderEvent represents an order event for Kafka
type OrderEvent struct {
	EventType   string      `json:"eventType"`
	OrderID     string      `json:"orderId"`
	UserID      string      `json:"userId"`
	Symbol      string      `json:"symbol"`
	Side        string      `json:"side"`      // BUY, SELL
	OrderType   string      `json:"orderType"` // MARKET, LIMIT, STOP
	Quantity    float64     `json:"quantity"`
	Price       float64     `json:"price,omitempty"`
	Status      string      `json:"status"`
	Timestamp   int64       `json:"timestamp"`
	PortfolioID string      `json:"portfolioId,omitempty"`
	AccountID   string      `json:"accountId,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// PublishOrderCreated publishes an order created event
func PublishOrderCreated(order *OrderEvent) error {
	order.EventType = EventOrderCreated
	order.Timestamp = time.Now().UnixMilli()
	return publishOrderEvent(order)
}

// PublishOrderCancelled publishes an order cancelled event
func PublishOrderCancelled(orderID, userID, reason string) error {
	event := &OrderEvent{
		EventType: EventOrderCancelled,
		OrderID:   orderID,
		UserID:    userID,
		Timestamp: time.Now().UnixMilli(),
		Metadata:  map[string]string{"reason": reason},
	}
	return publishOrderEvent(event)
}

// PublishOrderFilled publishes an order filled event
func PublishOrderFilled(orderID, userID string, filledQty, filledPrice float64) error {
	event := &OrderEvent{
		EventType: EventOrderFilled,
		OrderID:   orderID,
		UserID:    userID,
		Quantity:  filledQty,
		Price:     filledPrice,
		Timestamp: time.Now().UnixMilli(),
	}
	return publishOrderEvent(event)
}

// PublishOrderPartiallyFilled publishes a partially filled event
func PublishOrderPartiallyFilled(orderID, userID string, filledQty, remainingQty float64) error {
	event := &OrderEvent{
		EventType: EventOrderPartial,
		OrderID:   orderID,
		UserID:    userID,
		Quantity:  filledQty,
		Timestamp: time.Now().UnixMilli(),
		Metadata: map[string]float64{
			"remainingQty": remainingQty,
		},
	}
	return publishOrderEvent(event)
}

func publishOrderEvent(event *OrderEvent) error {
	if !IsConnected() {
		log.Printf("[Kafka] Not connected, skipping order event: %s", event.EventType)
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Printf("[Kafka] Publishing order event: %s for order %s", event.EventType, event.OrderID)
	return Publish(TopicOrders, event.OrderID, data)
}

// StartOrderConsumer starts consuming order events
func StartOrderConsumer(handler func(event *OrderEvent) error) {
	go func() {
		ConsumeOrders(func(key, value []byte) error {
			var event OrderEvent
			if err := json.Unmarshal(value, &event); err != nil {
				log.Printf("[Kafka] Failed to unmarshal order event: %v", err)
				return err
			}
			return handler(&event)
		})
	}()
}
