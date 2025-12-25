package cache

import (
	"fmt"
	"time"
)

// Quote represents a cached stock quote
type Quote struct {
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	PrevClose     float64 `json:"prevClose"`
	Volume        int64   `json:"volume"`
	Timestamp     int64   `json:"timestamp"`
}

const (
	// QuoteTTL is how long quotes are cached (5 seconds for real-time feel)
	QuoteTTL = 5 * time.Second

	// QuotePrefix is the Redis key prefix for quotes
	QuotePrefix = "quote:"
)

// SetQuote caches a stock quote
func SetQuote(symbol string, quote *Quote) error {
	key := QuotePrefix + symbol
	return SetJSON(key, quote, QuoteTTL)
}

// GetQuote retrieves a cached stock quote
func GetQuote(symbol string) (*Quote, error) {
	key := QuotePrefix + symbol
	quote := &Quote{}
	err := GetJSON(key, quote)
	if err != nil {
		return nil, err
	}
	return quote, nil
}

// GetQuotes retrieves multiple cached quotes
func GetQuotes(symbols []string) (map[string]*Quote, error) {
	result := make(map[string]*Quote)
	for _, symbol := range symbols {
		quote, err := GetQuote(symbol)
		if err == nil && quote != nil {
			result[symbol] = quote
		}
	}
	return result, nil
}

// DeleteQuote removes a cached quote
func DeleteQuote(symbol string) error {
	key := QuotePrefix + symbol
	return Delete(key)
}

// QuoteExists checks if a quote is cached
func QuoteExists(symbol string) (bool, error) {
	key := QuotePrefix + symbol
	return Exists(key)
}

// QuoteTTLRemaining returns remaining cache time for a quote
func QuoteTTLRemaining(symbol string) (time.Duration, error) {
	key := QuotePrefix + symbol
	return TTL(key)
}

// SetQuoteFromPrice creates a quote from price stream data
func SetQuoteFromPrice(symbol string, price float64, change float64, changePercent float64, volume int64) error {
	quote := &Quote{
		Symbol:        symbol,
		Price:         price,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        volume,
		Timestamp:     time.Now().UnixMilli(),
	}
	return SetQuote(symbol, quote)
}

// GetOrFetchQuote gets from cache or calls fetcher function
func GetOrFetchQuote(symbol string, fetcher func(string) (*Quote, error)) (*Quote, error) {
	// Try cache first
	quote, err := GetQuote(symbol)
	if err == nil && quote != nil {
		return quote, nil
	}

	// Fetch from source
	quote, err = fetcher(symbol)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if quote != nil {
		_ = SetQuote(symbol, quote)
	}

	return quote, nil
}

// CacheStats returns stats about cached quotes
func CacheStats() map[string]interface{} {
	keys, _ := Keys(QuotePrefix + "*")
	return map[string]interface{}{
		"cachedQuotes": len(keys),
		"connected":    IsConnected(),
	}
}

// InvalidateAllQuotes removes all cached quotes
func InvalidateAllQuotes() error {
	keys, err := Keys(QuotePrefix + "*")
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := Delete(key); err != nil {
			return fmt.Errorf("failed to delete %s: %w", key, err)
		}
	}
	return nil
}
