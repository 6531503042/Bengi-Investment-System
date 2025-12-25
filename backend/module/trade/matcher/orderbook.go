package matcher

import (
	"sort"
	"sync"
)

// Order represents an order in the order book
type Order struct {
	ID          string
	UserID      string
	Symbol      string
	Side        string  // BUY, SELL
	Type        string  // LIMIT, MARKET
	Price       float64 // 0 for MARKET orders
	Quantity    float64
	FilledQty   float64
	Timestamp   int64
	PortfolioID string
	AccountID   string
}

// OrderBook manages buy and sell orders for a symbol
type OrderBook struct {
	Symbol     string
	BuyOrders  []*Order // Sorted by price DESC, then time ASC
	SellOrders []*Order // Sorted by price ASC, then time ASC
	mu         sync.RWMutex
}

// NewOrderBook creates a new order book for a symbol
func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		Symbol:     symbol,
		BuyOrders:  make([]*Order, 0),
		SellOrders: make([]*Order, 0),
	}
}

// AddOrder adds an order to the book
func (ob *OrderBook) AddOrder(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Side == "BUY" {
		ob.BuyOrders = append(ob.BuyOrders, order)
		ob.sortBuyOrders()
	} else {
		ob.SellOrders = append(ob.SellOrders, order)
		ob.sortSellOrders()
	}
}

// RemoveOrder removes an order from the book
func (ob *OrderBook) RemoveOrder(orderID string) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Check buy orders
	for i, o := range ob.BuyOrders {
		if o.ID == orderID {
			ob.BuyOrders = append(ob.BuyOrders[:i], ob.BuyOrders[i+1:]...)
			return true
		}
	}

	// Check sell orders
	for i, o := range ob.SellOrders {
		if o.ID == orderID {
			ob.SellOrders = append(ob.SellOrders[:i], ob.SellOrders[i+1:]...)
			return true
		}
	}

	return false
}

// GetBestBid returns the highest buy price
func (ob *OrderBook) GetBestBid() float64 {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if len(ob.BuyOrders) == 0 {
		return 0
	}
	return ob.BuyOrders[0].Price
}

// GetBestAsk returns the lowest sell price
func (ob *OrderBook) GetBestAsk() float64 {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if len(ob.SellOrders) == 0 {
		return 0
	}
	return ob.SellOrders[0].Price
}

// GetSpread returns the bid-ask spread
func (ob *OrderBook) GetSpread() float64 {
	return ob.GetBestAsk() - ob.GetBestBid()
}

// GetDepth returns the order book depth
func (ob *OrderBook) GetDepth() (bids, asks int) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	return len(ob.BuyOrders), len(ob.SellOrders)
}

// Sort buy orders: highest price first, then oldest first
func (ob *OrderBook) sortBuyOrders() {
	sort.Slice(ob.BuyOrders, func(i, j int) bool {
		if ob.BuyOrders[i].Price != ob.BuyOrders[j].Price {
			return ob.BuyOrders[i].Price > ob.BuyOrders[j].Price
		}
		return ob.BuyOrders[i].Timestamp < ob.BuyOrders[j].Timestamp
	})
}

// Sort sell orders: lowest price first, then oldest first
func (ob *OrderBook) sortSellOrders() {
	sort.Slice(ob.SellOrders, func(i, j int) bool {
		if ob.SellOrders[i].Price != ob.SellOrders[j].Price {
			return ob.SellOrders[i].Price < ob.SellOrders[j].Price
		}
		return ob.SellOrders[i].Timestamp < ob.SellOrders[j].Timestamp
	})
}
