package matcher

import (
	"log"
	"sync"
	"time"
)

// Match represents a successful order match
type Match struct {
	BuyOrderID  string
	SellOrderID string
	Symbol      string
	Price       float64
	Quantity    float64
	BuyerID     string
	SellerID    string
	Timestamp   int64
}

// MatchHandler is called when orders are matched
type MatchHandler func(match *Match) error

// Engine is the order matching engine
type Engine struct {
	orderBooks   map[string]*OrderBook
	mu           sync.RWMutex
	matchHandler MatchHandler
	running      bool
	stopCh       chan struct{}
}

// NewEngine creates a new matching engine
func NewEngine(handler MatchHandler) *Engine {
	return &Engine{
		orderBooks:   make(map[string]*OrderBook),
		matchHandler: handler,
		stopCh:       make(chan struct{}),
	}
}

// Start starts the matching engine
func (e *Engine) Start() {
	e.running = true
	log.Println("[Matcher] Engine started")

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-e.stopCh:
				return
			case <-ticker.C:
				e.matchAllBooks()
			}
		}
	}()
}

// Stop stops the matching engine
func (e *Engine) Stop() {
	e.running = false
	close(e.stopCh)
	log.Println("[Matcher] Engine stopped")
}

// AddOrder adds an order to the appropriate order book
func (e *Engine) AddOrder(order *Order) {
	e.mu.Lock()
	defer e.mu.Unlock()

	book, exists := e.orderBooks[order.Symbol]
	if !exists {
		book = NewOrderBook(order.Symbol)
		e.orderBooks[order.Symbol] = book
	}

	book.AddOrder(order)
	log.Printf("[Matcher] Added %s order for %s: %.2f @ %.2f", order.Side, order.Symbol, order.Quantity, order.Price)
}

// CancelOrder removes an order from the book
func (e *Engine) CancelOrder(symbol, orderID string) bool {
	e.mu.RLock()
	book, exists := e.orderBooks[symbol]
	e.mu.RUnlock()

	if !exists {
		return false
	}

	return book.RemoveOrder(orderID)
}

// matchAllBooks runs matching on all order books
func (e *Engine) matchAllBooks() {
	e.mu.RLock()
	symbols := make([]string, 0, len(e.orderBooks))
	for symbol := range e.orderBooks {
		symbols = append(symbols, symbol)
	}
	e.mu.RUnlock()

	for _, symbol := range symbols {
		e.matchBook(symbol)
	}
}

// matchBook matches orders in a single order book
func (e *Engine) matchBook(symbol string) {
	e.mu.RLock()
	book, exists := e.orderBooks[symbol]
	e.mu.RUnlock()

	if !exists {
		return
	}

	book.mu.Lock()
	defer book.mu.Unlock()

	for len(book.BuyOrders) > 0 && len(book.SellOrders) > 0 {
		buy := book.BuyOrders[0]
		sell := book.SellOrders[0]

		// Check if orders can match
		// For MARKET orders, always match
		// For LIMIT orders, buy price must >= sell price
		canMatch := false
		matchPrice := 0.0

		if buy.Type == "MARKET" || sell.Type == "MARKET" {
			canMatch = true
			if sell.Type == "MARKET" {
				matchPrice = buy.Price
			} else {
				matchPrice = sell.Price
			}
		} else if buy.Price >= sell.Price {
			canMatch = true
			// Use the earlier order's price (price-time priority)
			if buy.Timestamp < sell.Timestamp {
				matchPrice = buy.Price
			} else {
				matchPrice = sell.Price
			}
		}

		if !canMatch {
			break
		}

		// Calculate match quantity
		matchQty := min(buy.Quantity-buy.FilledQty, sell.Quantity-sell.FilledQty)

		// Create match
		match := &Match{
			BuyOrderID:  buy.ID,
			SellOrderID: sell.ID,
			Symbol:      symbol,
			Price:       matchPrice,
			Quantity:    matchQty,
			BuyerID:     buy.UserID,
			SellerID:    sell.UserID,
			Timestamp:   time.Now().UnixMilli(),
		}

		// Update filled quantities
		buy.FilledQty += matchQty
		sell.FilledQty += matchQty

		// Remove fully filled orders
		if buy.FilledQty >= buy.Quantity {
			book.BuyOrders = book.BuyOrders[1:]
		}
		if sell.FilledQty >= sell.Quantity {
			book.SellOrders = book.SellOrders[1:]
		}

		// Call match handler
		if e.matchHandler != nil {
			if err := e.matchHandler(match); err != nil {
				log.Printf("[Matcher] Error handling match: %v", err)
			}
		}

		log.Printf("[Matcher] Matched: %s %.4f @ %.2f", symbol, matchQty, matchPrice)
	}
}

// GetOrderBookStats returns stats for an order book
func (e *Engine) GetOrderBookStats(symbol string) map[string]interface{} {
	e.mu.RLock()
	book, exists := e.orderBooks[symbol]
	e.mu.RUnlock()

	if !exists {
		return nil
	}

	bids, asks := book.GetDepth()
	return map[string]interface{}{
		"symbol":   symbol,
		"bestBid":  book.GetBestBid(),
		"bestAsk":  book.GetBestAsk(),
		"spread":   book.GetSpread(),
		"bidDepth": bids,
		"askDepth": asks,
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
