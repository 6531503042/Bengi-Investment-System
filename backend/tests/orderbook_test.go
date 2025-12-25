package tests

import (
	"testing"

	"github.com/bricksocoolxd/bengi-investment-system/module/trade/matcher"
)

func TestOrderBook_AddOrder(t *testing.T) {
	book := matcher.NewOrderBook("AAPL")

	order := &matcher.Order{
		ID:        "order-1",
		UserID:    "user-1",
		Symbol:    "AAPL",
		Side:      "BUY",
		Type:      "LIMIT",
		Price:     150.00,
		Quantity:  10,
		Timestamp: 1000,
	}

	book.AddOrder(order)

	if book.GetBestBid() != 150.00 {
		t.Errorf("Expected best bid 150.00, got %f", book.GetBestBid())
	}
}

func TestOrderBook_BidAskSpread(t *testing.T) {
	book := matcher.NewOrderBook("AAPL")

	buyOrder := &matcher.Order{
		ID:        "buy-1",
		Side:      "BUY",
		Type:      "LIMIT",
		Price:     149.00,
		Quantity:  10,
		Timestamp: 1000,
	}
	book.AddOrder(buyOrder)

	sellOrder := &matcher.Order{
		ID:        "sell-1",
		Side:      "SELL",
		Type:      "LIMIT",
		Price:     151.00,
		Quantity:  10,
		Timestamp: 1001,
	}
	book.AddOrder(sellOrder)

	if book.GetBestBid() != 149.00 {
		t.Errorf("Expected best bid 149.00, got %f", book.GetBestBid())
	}

	if book.GetBestAsk() != 151.00 {
		t.Errorf("Expected best ask 151.00, got %f", book.GetBestAsk())
	}

	if book.GetSpread() != 2.00 {
		t.Errorf("Expected spread 2.00, got %f", book.GetSpread())
	}
}

func TestOrderBook_RemoveOrder(t *testing.T) {
	book := matcher.NewOrderBook("AAPL")

	order := &matcher.Order{
		ID:       "order-1",
		Side:     "BUY",
		Price:    150.00,
		Quantity: 10,
	}
	book.AddOrder(order)

	removed := book.RemoveOrder("order-1")
	if !removed {
		t.Error("Expected order to be removed")
	}

	bids, _ := book.GetDepth()
	if bids != 0 {
		t.Errorf("Expected 0 bids, got %d", bids)
	}
}

func TestOrderBook_PriceTimePriority(t *testing.T) {
	book := matcher.NewOrderBook("AAPL")

	// Add orders with same price, different times
	order1 := &matcher.Order{ID: "order-1", Side: "BUY", Price: 150.00, Quantity: 10, Timestamp: 1000}
	order2 := &matcher.Order{ID: "order-2", Side: "BUY", Price: 150.00, Quantity: 10, Timestamp: 999}  // Earlier
	order3 := &matcher.Order{ID: "order-3", Side: "BUY", Price: 151.00, Quantity: 10, Timestamp: 1001} // Higher price

	book.AddOrder(order1)
	book.AddOrder(order2)
	book.AddOrder(order3)

	// Best bid should be highest price
	if book.GetBestBid() != 151.00 {
		t.Errorf("Expected best bid 151.00 (highest price), got %f", book.GetBestBid())
	}
}
