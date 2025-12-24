package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"github.com/gorilla/websocket"
)

// Finnhub message types
type (
	// FinnhubTradeMessage represents incoming trade data from Finnhub
	FinnhubMessage struct {
		Type string         `json:"type"`
		Data []FinnhubTrade `json:"data,omitempty"`
	}

	FinnhubTrade struct {
		Symbol    string   `json:"s"`
		Price     float64  `json:"p"`
		Volume    float64  `json:"v"`
		Timestamp int64    `json:"t"`
		Condition []string `json:"c,omitempty"`
	}

	// FinnhubSubscribe is the subscribe message format
	FinnhubSubscribe struct {
		Type   string `json:"type"`
		Symbol string `json:"symbol"`
	}

	// PriceStream manages connection to Finnhub WebSocket
	PriceStream struct {
		conn           *websocket.Conn
		apiKey         string
		symbols        map[string]bool
		lastPrices     map[string]*PricePayload // Track last prices for change calculation
		mu             sync.RWMutex
		done           chan struct{}
		reconnectDelay time.Duration
		isConnected    bool
	}
)

var (
	priceStream *PriceStream
	streamOnce  sync.Once
)

const (
	finnhubWSURL      = "wss://ws.finnhub.io"
	maxReconnectDelay = 30 * time.Second
)

// GetPriceStream returns singleton price stream instance
func GetPriceStream() *PriceStream {
	streamOnce.Do(func() {
		priceStream = &PriceStream{
			apiKey:         config.AppConfig.FinnhubAPIKey,
			symbols:        make(map[string]bool),
			lastPrices:     make(map[string]*PricePayload),
			done:           make(chan struct{}),
			reconnectDelay: time.Second,
		}
	})
	return priceStream
}

// Start starts the price stream connection
func (ps *PriceStream) Start() error {
	if ps.apiKey == "" {
		log.Println("[PriceStream] No Finnhub API key configured, skipping")
		return nil
	}
	go ps.connectLoop()
	return nil
}

// Stop stops the price stream
func (ps *PriceStream) Stop() {
	close(ps.done)
	if ps.conn != nil {
		ps.conn.Close()
	}
}

// connectLoop handles connection and reconnection
func (ps *PriceStream) connectLoop() {
	for {
		select {
		case <-ps.done:
			return
		default:
			if err := ps.connect(); err != nil {
				log.Printf("[PriceStream] Connection error: %v", err)
				ps.isConnected = false

				// Exponential backoff
				time.Sleep(ps.reconnectDelay)
				ps.reconnectDelay *= 2
				if ps.reconnectDelay > maxReconnectDelay {
					ps.reconnectDelay = maxReconnectDelay
				}
				continue
			}

			// Reset reconnect delay on successful connect
			ps.reconnectDelay = time.Second
			ps.isConnected = true

			// Read messages
			ps.readLoop()
		}
	}
}

// connect establishes WebSocket connection to Finnhub
func (ps *PriceStream) connect() error {
	url := finnhubWSURL + "?token=" + ps.apiKey

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	ps.conn = conn
	log.Println("[PriceStream] Connected to Finnhub")

	// Resubscribe to existing symbols
	ps.mu.RLock()
	symbols := make([]string, 0, len(ps.symbols))
	for sym := range ps.symbols {
		symbols = append(symbols, sym)
	}
	ps.mu.RUnlock()

	for _, sym := range symbols {
		ps.sendSubscribe(sym)
	}

	return nil
}

// readLoop reads messages from Finnhub
func (ps *PriceStream) readLoop() {
	defer func() {
		if ps.conn != nil {
			ps.conn.Close()
		}
	}()

	for {
		select {
		case <-ps.done:
			return
		default:
			_, message, err := ps.conn.ReadMessage()
			if err != nil {
				log.Printf("[PriceStream] Read error: %v", err)
				return
			}

			ps.handleMessage(message)
		}
	}
}

// handleMessage processes incoming messages from Finnhub
func (ps *PriceStream) handleMessage(data []byte) {
	var msg FinnhubMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[PriceStream] Parse error: %v", err)
		return
	}

	switch msg.Type {
	case "trade":
		ps.handleTrades(msg.Data)
	case "ping":
		// Respond to ping with pong
		ps.conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"pong"}`))
	case "error":
		log.Printf("[PriceStream] Error from Finnhub: %s", string(data))
	default:
		// Ignore other message types
	}
}

// handleTrades processes trade data and broadcasts to subscribers
func (ps *PriceStream) handleTrades(trades []FinnhubTrade) {
	// Group trades by symbol and use latest price
	symbolPrices := make(map[string]*FinnhubTrade)
	for i := range trades {
		trade := &trades[i]
		existing, ok := symbolPrices[trade.Symbol]
		if !ok || trade.Timestamp > existing.Timestamp {
			symbolPrices[trade.Symbol] = trade
		}
	}

	// Broadcast each symbol's latest price
	for symbol, trade := range symbolPrices {
		ps.mu.Lock()
		lastPrice := ps.lastPrices[symbol]

		var change, changePercent float64
		if lastPrice != nil && lastPrice.Price > 0 {
			change = trade.Price - lastPrice.Price
			changePercent = (change / lastPrice.Price) * 100
		}

		payload := &PricePayload{
			Symbol:        symbol,
			Price:         trade.Price,
			Change:        change,
			ChangePercent: changePercent,
			Volume:        int64(trade.Volume),
		}

		ps.lastPrices[symbol] = payload
		ps.mu.Unlock()

		// Publish to Event Bus
		PublishPrice(symbol, payload)
	}
}

// Subscribe adds symbols to watch
func (ps *PriceStream) Subscribe(symbols ...string) {
	ps.mu.Lock()
	newSymbols := make([]string, 0)
	for _, sym := range symbols {
		if !ps.symbols[sym] {
			ps.symbols[sym] = true
			newSymbols = append(newSymbols, sym)
		}
	}
	ps.mu.Unlock()

	if ps.isConnected {
		for _, sym := range newSymbols {
			ps.sendSubscribe(sym)
		}
	}
}

// Unsubscribe removes symbols from watch
func (ps *PriceStream) Unsubscribe(symbols ...string) {
	ps.mu.Lock()
	removeSymbols := make([]string, 0)
	for _, sym := range symbols {
		if ps.symbols[sym] {
			delete(ps.symbols, sym)
			removeSymbols = append(removeSymbols, sym)
		}
	}
	ps.mu.Unlock()

	if ps.isConnected {
		for _, sym := range removeSymbols {
			ps.sendUnsubscribe(sym)
		}
	}
}

// sendSubscribe sends subscribe message to Finnhub
func (ps *PriceStream) sendSubscribe(symbol string) {
	msg := FinnhubSubscribe{
		Type:   "subscribe",
		Symbol: symbol,
	}

	data, _ := json.Marshal(msg)
	if err := ps.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[PriceStream] Subscribe error: %v", err)
	} else {
		log.Printf("[PriceStream] Subscribed to: %s", symbol)
	}
}

// sendUnsubscribe sends unsubscribe message to Finnhub
func (ps *PriceStream) sendUnsubscribe(symbol string) {
	msg := FinnhubSubscribe{
		Type:   "unsubscribe",
		Symbol: symbol,
	}

	data, _ := json.Marshal(msg)
	if err := ps.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[PriceStream] Unsubscribe error: %v", err)
	} else {
		log.Printf("[PriceStream] Unsubscribed from: %s", symbol)
	}
}

// GetSubscribedSymbols returns currently subscribed symbols
func (ps *PriceStream) GetSubscribedSymbols() []string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	symbols := make([]string, 0, len(ps.symbols))
	for sym := range ps.symbols {
		symbols = append(symbols, sym)
	}
	return symbols
}

// IsConnected returns connection status
func (ps *PriceStream) IsConnected() bool {
	return ps.isConnected
}

// GetLastPrice returns the last known price for a symbol
func (ps *PriceStream) GetLastPrice(symbol string) *PricePayload {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.lastPrices[symbol]
}
