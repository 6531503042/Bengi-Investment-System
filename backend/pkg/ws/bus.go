package ws

import "sync"

type Subscriber func(msg *Message)

type EventBus struct {
	subscribers map[string]map[string]Subscriber
	mu          sync.RWMutex
}

var Bus *EventBus

func InitBus() {
	Bus = NewEventBus()
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string]map[string]Subscriber),
	}
}

func (eb *EventBus) Subscribe(topic, subscriberID string, callback Subscriber) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.subscribers[topic] == nil {
		eb.subscribers[topic] = make(map[string]Subscriber)
	}
	eb.subscribers[topic][subscriberID] = callback
}

func (eb *EventBus) Unsubscribe(topic, subscribeID string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.subscribers[topic] != nil {
		delete(eb.subscribers[topic], subscribeID)
		if len(eb.subscribers[topic]) == 0 {
			delete(eb.subscribers, topic)
		}
	}
}

func (eb *EventBus) UnsubscribeAll(subscriberID string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for topic := range eb.subscribers {
		delete(eb.subscribers[topic], subscriberID)
		if len(eb.subscribers[topic]) == 0 {
			delete(eb.subscribers, topic)
		}
	}
}

func (eb *EventBus) Publish(topic string, msg *Message) {
	eb.mu.RLock()
	subs := eb.subscribers[topic]
	eb.mu.Unlock()

	for _, callback := range subs {
		go callback(msg)
	}
}

func (eb *EventBus) PublishBytes(topic string, data []byte) {
	msg, err := ParseMessage(data)
	if err != nil {
		return
	}
	eb.Publish(topic, msg)
}

func (eb *EventBus) HasSubscribers(topic string) bool {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	//If subscriber map is nil, it means there are no subscribers for this topic
	return len(eb.subscribers[topic]) > 0
}

func (eb *EventBus) GetTopics() []string {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	topics := make([]string, 0, len(eb.subscribers))
	for topic := range eb.subscribers {
		topics = append(topics, topic)
	}
	return topics
}

func (eb *EventBus) SubscriberCount(topic string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.subscribers[topic])
}

// ========== Helper Functions for Publishing ==========
// PublishPrice publishes a price update
func PublishPrice(symbol string, payload *PricePayload) {
	topic := TopicPrice(symbol)
	msg := NewMessage(TypePriceUpdate, topic, payload)
	Bus.Publish(topic, msg)
}

// PublishOrderUpdate publishes an order update to user
func PublishOrderUpdate(userID string, payload *OrderPayload) {
	topic := TopicOrder(userID)
	msg := NewMessage(TypeOrderUpdate, topic, payload)
	Bus.Publish(topic, msg)
}

// PublishTradeUpdate publishes a trade update to user
func PublishTradeUpdate(userID string, payload *TradePayload) {
	topic := TopicTrade(userID)
	msg := NewMessage(TypeTradeUpdate, topic, payload)
	Bus.Publish(topic, msg)
}
