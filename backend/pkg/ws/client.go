package ws

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// Client represents a WebSocket connection
type Client struct {
	ID            string
	UserID        string // empty if not authenticated
	Conn          *websocket.Conn
	Send          chan []byte
	subscriptions map[string]bool
	mu            sync.RWMutex
	closed        bool
}

// ClientManager manages all connected clients
type ClientManager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

// Global client manager
var Manager *ClientManager

// InitManager initializes the client manager
func InitManager() {
	Manager = &ClientManager{
		clients: make(map[string]*Client),
	}
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, userID string) *Client {
	return &Client{
		ID:            uuid.New().String(),
		UserID:        userID,
		Conn:          conn,
		Send:          make(chan []byte, 256),
		subscriptions: make(map[string]bool),
	}
}

// Register adds a client to the manager
func (m *ClientManager) Register(client *Client) {
	m.mu.Lock()
	m.clients[client.ID] = client
	m.mu.Unlock()
	log.Printf("[WS] Client connected: %s (User: %s)", client.ID, client.UserID)
}

// Unregister removes a client from the manager
func (m *ClientManager) Unregister(client *Client) {
	m.mu.Lock()
	delete(m.clients, client.ID)
	m.mu.Unlock()

	// Unsubscribe from all topics
	Bus.UnsubscribeAll(client.ID)
	log.Printf("[WS] Client disconnected: %s", client.ID)
}

// GetClient gets a client by ID
func (m *ClientManager) GetClient(id string) *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[id]
}

// ClientCount returns the number of connected clients
func (m *ClientManager) ClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

// ========== Client Methods ==========

// Subscribe subscribes client to a topic
func (c *Client) Subscribe(topic string) {
	// Validate topic
	if !ValidateTopic(topic) {
		c.SendError("INVALID_TOPIC", "Invalid topic format")
		return
	}

	// Check authorization for user topics
	if IsUserTopic(topic) {
		topicUser := GetUserFromTopic(topic)
		if c.UserID == "" || c.UserID != topicUser {
			c.SendError("UNAUTHORIZED", "Cannot subscribe to this topic")
			return
		}
	}

	c.mu.Lock()
	c.subscriptions[topic] = true
	c.mu.Unlock()

	// Subscribe to event bus
	Bus.Subscribe(topic, c.ID, func(msg *Message) {
		c.Send <- msg.ToBytes()
	})

	if len(topic) > len(TopicPricePrefix) && topic[:len(TopicPricePrefix)] == TopicPricePrefix {
		symbol := topic[len(TopicPricePrefix):]
		GetPriceStream().Subscribe(symbol)
	}

	// Confirm subscription
	c.Send <- NewMessage(TypeSubscribed, topic, nil).ToBytes()
	log.Printf("[WS] Client %s subscribed to %s", c.ID, topic)
}

// Unsubscribe unsubscribes client from a topic
func (c *Client) Unsubscribe(topic string) {
	c.mu.Lock()
	delete(c.subscriptions, topic)
	c.mu.Unlock()

	Bus.Unsubscribe(topic, c.ID)
	c.Send <- NewMessage(TypeUnsubscribed, topic, nil).ToBytes()
	log.Printf("[WS] Client %s unsubscribed from %s", c.ID, topic)
}

// SendError sends an error message to client
func (c *Client) SendError(code, message string) {
	msg := NewMessage(TypeError, "", &ErrorPayload{
		Code:    code,
		Message: message,
	})
	c.Send <- msg.ToBytes()
}

// Close closes the client connection
func (c *Client) Close() {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true
	c.mu.Unlock()

	close(c.Send)
}

// ReadPump reads messages from WebSocket
func (c *Client) ReadPump() {
	defer func() {
		Manager.Unregister(c)
		c.Close()
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		msg, err := ParseMessage(data)
		if err != nil {
			c.SendError("PARSE_ERROR", "Invalid message format")
			continue
		}

		c.handleMessage(msg)
	}
}

// WritePump writes messages to WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming messages
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case TypeSubscribe:
		c.Subscribe(msg.Topic)
	case TypeUnsubscribe:
		c.Unsubscribe(msg.Topic)
	case TypePing:
		c.Send <- NewMessage(TypePong, "", nil).ToBytes()
	default:
		c.SendError("UNKNOWN_TYPE", "Unknown message type")
	}
}
