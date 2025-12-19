package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const InstrumentCollection = "instruments"

type InstrumentType string

const (
	InstrumentTypeSock   InstrumentType = "Stock"
	InstrumentTypeETF    InstrumentType = "ETF"
	InstrumentTypeCrypto InstrumentType = "Crypto"
	InstrumentTypeFuture InstrumentType = "Future"
	InstrumentTypeOption InstrumentType = "Option"
)

type InstrumentStatus string

const (
	InstrumentStatusActive   InstrumentStatus = "Active"
	InstrumentStatusInactive InstrumentStatus = "Inactive"
	InstrumentStatusDelisted InstrumentStatus = "DELISTED"
)

type (
	Instrument struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Symbol      string             `bson:"symbol" json:"symbol"`     // AAPL, BTC/USD
		Name        string             `bson:"name" json:"name"`         // Apple Inc.
		Type        InstrumentType     `bson:"type" json:"type"`         // STOCK, ETF, CRYPTO
		Exchange    string             `bson:"exchange" json:"exchange"` // NASDAQ, NYSE
		Currency    string             `bson:"currency" json:"currency"` // USD, THB
		Status      InstrumentStatus   `bson:"status" json:"status"`
		Description string             `bson:"description" json:"description"`
		LogoURL     string             `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`
		CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
		UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	}

	Quote struct {
		Symbol        string    `json:"symbol"`
		Price         float64   `json:"price"`
		Open          float64   `json:"open"`
		High          float64   `json:"high"`
		Low           float64   `json:"low"`
		Close         float64   `json:"close"`
		PreviousClose float64   `json:"previousClose"`
		Volume        int64     `json:"volume"`
		Change        float64   `json:"change"`
		ChangePercent float64   `json:"changePercent"`
		Timestamp     time.Time `json:"timestamp"`
	}
)
