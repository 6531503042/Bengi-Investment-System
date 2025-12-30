package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InstrumentCollection is the MongoDB collection name for tradeable instruments.
const InstrumentCollection = "instruments"

// InstrumentType categorizes the type of financial instrument.
type InstrumentType string

const (
	InstrumentTypeStock  InstrumentType = "Stock"  // Individual company shares
	InstrumentTypeETF    InstrumentType = "ETF"    // Exchange-traded funds
	InstrumentTypeCrypto InstrumentType = "Crypto" // Cryptocurrencies
	InstrumentTypeFuture InstrumentType = "Future" // Futures contracts
	InstrumentTypeOption InstrumentType = "Option" // Options contracts
)

// InstrumentStatus indicates whether an instrument is available for trading.
type InstrumentStatus string

const (
	InstrumentStatusActive   InstrumentStatus = "Active"   // Available for trading
	InstrumentStatusInactive InstrumentStatus = "Inactive" // Temporarily unavailable
	InstrumentStatusDelisted InstrumentStatus = "Delisted" // Permanently removed
)

// Instrument represents a tradeable financial instrument.
// Can be stocks, ETFs, cryptocurrencies, futures, or options.
type Instrument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Symbol      string             `bson:"symbol" json:"symbol"`           // AAPL, BTC/USD, etc.
	Name        string             `bson:"name" json:"name"`               // Apple Inc.
	Type        InstrumentType     `bson:"type" json:"type"`               // Stock, ETF, Crypto
	Exchange    string             `bson:"exchange" json:"exchange"`       // NASDAQ, NYSE, Binance
	Currency    string             `bson:"currency" json:"currency"`       // Base currency (USD, THB)
	Status      InstrumentStatus   `bson:"status" json:"status"`           // Active, Inactive, Delisted
	Description string             `bson:"description" json:"description"` // Brief info about the instrument
	LogoURL     string             `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Quote represents real-time price data for an instrument.
// Updated via market data feeds (Twelve Data, Yahoo Finance).
type Quote struct {
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`         // Current/last traded price
	Open          float64   `json:"open"`          // Opening price for the day
	High          float64   `json:"high"`          // Day's high
	Low           float64   `json:"low"`           // Day's low
	Close         float64   `json:"close"`         // Previous close (or current if market closed)
	PreviousClose float64   `json:"previousClose"` // Previous trading day's close
	Volume        int64     `json:"volume"`        // Trading volume
	Change        float64   `json:"change"`        // Price change from previous close
	ChangePercent float64   `json:"changePercent"` // Percentage change
	Timestamp     time.Time `json:"timestamp"`     // Time of last update
}

// IsStock returns true if this instrument is a stock.
func (i *Instrument) IsStock() bool {
	return i.Type == InstrumentTypeStock
}

// IsCrypto returns true if this is a cryptocurrency.
func (i *Instrument) IsCrypto() bool {
	return i.Type == InstrumentTypeCrypto
}
