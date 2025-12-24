package ws

import (
	"encoding/json"
	"time"
)

const (
	TypeSubscribe    = "SUBSCRIBE"
	TypeUnsubscribe  = "UNSUBSCRIBE"
	TypePing         = "PING"
	TypePong         = "PONG"
	TypeSubscribed   = "SUBSCRIBED"
	TypeUnsubscribed = "UNSUBSCRIBED"
	TypePriceUpdate  = "PRICE_UPDATE"
	TypeOrderUpdate  = "ORDER_UPDATE"
	TypeTradeUpdate  = "TRADE_UPDATE"
	TypeError        = "ERROR"
)

type Message struct {
	Type      string          `json:"type"`
	Topic     string          `json:"topic,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
}

// PricePayload for price updates
type PricePayload struct {
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"price"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	Volume        int64   `json:"volume"`
}

// OrderPayload for order updates
type OrderPayload struct {
	OrderID   string  `json:"orderId"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Status    string  `json:"status"`
	FilledQty float64 `json:"filledQty"`
	AvgPrice  float64 `json:"avgPrice,omitempty"`
}

// TradePayload for trade updates
type TradePayload struct {
	TradeID    string  `json:"tradeId"`
	OrderID    string  `json:"orderId"`
	Symbol     string  `json:"symbol"`
	Side       string  `json:"side"`
	Quantity   float64 `json:"quantity"`
	Price      float64 `json:"price"`
	Commission float64 `json:"commission"`
}

// ErrorPayload for error messages
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewMessage(msgType, topic string, data interface{}) *Message {
	var rawData json.RawMessage
	if data != nil {
		rawData, _ = json.Marshal(data)
	}

	return &Message{
		Type:      msgType,
		Topic:     topic,
		Data:      rawData,
		Timestamp: time.Now().UnixMilli(),
	}
}

func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ToBytes serializes the Message to JSON bytes
func (m *Message) ToBytes() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("{}")
	}
	return data
}
