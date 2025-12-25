package tests

import (
	"testing"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/trade/matcher"
)

func TestMatchEngine_BasicMatch(t *testing.T) {
	matches := make([]*matcher.Match, 0)

	engine := matcher.NewEngine(func(match *matcher.Match) error {
		matches = append(matches, match)
		return nil
	})

	engine.Start()
	defer engine.Stop()

	// Add a buy order
	buyOrder := &matcher.Order{
		ID:        "buy-1",
		UserID:    "buyer-1",
		Symbol:    "AAPL",
		Side:      "BUY",
		Type:      "LIMIT",
		Price:     150.00,
		Quantity:  10,
		Timestamp: time.Now().UnixMilli(),
	}
	engine.AddOrder(buyOrder)

	// Add a matching sell order
	sellOrder := &matcher.Order{
		ID:        "sell-1",
		UserID:    "seller-1",
		Symbol:    "AAPL",
		Side:      "SELL",
		Type:      "LIMIT",
		Price:     149.00, // Lower than buy, should match
		Quantity:  10,
		Timestamp: time.Now().UnixMilli(),
	}
	engine.AddOrder(sellOrder)

	// Wait for matching
	time.Sleep(200 * time.Millisecond)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
		return
	}

	if matches[0].Quantity != 10 {
		t.Errorf("Expected match quantity 10, got %f", matches[0].Quantity)
	}
}

func TestMatchEngine_PartialMatch(t *testing.T) {
	matches := make([]*matcher.Match, 0)

	engine := matcher.NewEngine(func(match *matcher.Match) error {
		matches = append(matches, match)
		return nil
	})

	engine.Start()
	defer engine.Stop()

	// Buy order for 20 shares
	buyOrder := &matcher.Order{
		ID:        "buy-1",
		UserID:    "buyer-1",
		Symbol:    "AAPL",
		Side:      "BUY",
		Type:      "LIMIT",
		Price:     150.00,
		Quantity:  20,
		Timestamp: time.Now().UnixMilli(),
	}
	engine.AddOrder(buyOrder)

	// Sell order for only 10 shares
	sellOrder := &matcher.Order{
		ID:        "sell-1",
		UserID:    "seller-1",
		Symbol:    "AAPL",
		Side:      "SELL",
		Type:      "LIMIT",
		Price:     149.00,
		Quantity:  10,
		Timestamp: time.Now().UnixMilli(),
	}
	engine.AddOrder(sellOrder)

	time.Sleep(200 * time.Millisecond)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
		return
	}

	// Should match 10 (the smaller quantity)
	if matches[0].Quantity != 10 {
		t.Errorf("Expected match quantity 10, got %f", matches[0].Quantity)
	}
}

func TestMatchEngine_NoMatchWhenPricesDontCross(t *testing.T) {
	matches := make([]*matcher.Match, 0)

	engine := matcher.NewEngine(func(match *matcher.Match) error {
		matches = append(matches, match)
		return nil
	})

	engine.Start()
	defer engine.Stop()

	// Buy at 100
	buyOrder := &matcher.Order{
		ID:       "buy-1",
		Symbol:   "AAPL",
		Side:     "BUY",
		Type:     "LIMIT",
		Price:    100.00,
		Quantity: 10,
	}
	engine.AddOrder(buyOrder)

	// Sell at 110 (higher than buy, won't match)
	sellOrder := &matcher.Order{
		ID:       "sell-1",
		Symbol:   "AAPL",
		Side:     "SELL",
		Type:     "LIMIT",
		Price:    110.00,
		Quantity: 10,
	}
	engine.AddOrder(sellOrder)

	time.Sleep(200 * time.Millisecond)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches (prices don't cross), got %d", len(matches))
	}
}
