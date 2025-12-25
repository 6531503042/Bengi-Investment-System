package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"github.com/segmentio/kafka-go"
)

var (
	writer *kafka.Writer
	ctx    = context.Background()
)

// Topics
const (
	TopicOrders = "orders"
	TopicTrades = "trades"
	TopicEvents = "events"
)

// Initialize sets up Kafka writer
func Initialize() error {
	brokers := strings.Split(config.AppConfig.KafkaBrokers, ",")

	writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	}

	// Test connection by writing to a temp topic
	testConn, err := kafka.DialLeader(ctx, "tcp", brokers[0], "test", 0)
	if err != nil {
		log.Printf("[Kafka] Failed to connect: %v", err)
		return err
	}
	testConn.Close()

	log.Println("✅ Connected to Kafka")
	return nil
}

// IsConnected checks if Kafka is available
func IsConnected() bool {
	if writer == nil {
		return false
	}
	brokers := strings.Split(config.AppConfig.KafkaBrokers, ",")
	conn, err := kafka.DialLeader(ctx, "tcp", brokers[0], "test", 0)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Close closes the Kafka writer
func Close() error {
	if writer != nil {
		return writer.Close()
	}
	return nil
}

// Publish sends a message to a topic
func Publish(topic string, key string, value interface{}) error {
	if writer == nil {
		return nil // Fail silently if Kafka not initialized
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	})
}

// PublishOrder publishes an order event
func PublishOrder(orderID string, order interface{}) error {
	return Publish(TopicOrders, orderID, order)
}

// PublishTrade publishes a trade event
func PublishTrade(tradeID string, trade interface{}) error {
	return Publish(TopicTrades, tradeID, trade)
}

// PublishEvent publishes a generic event
func PublishEvent(eventType string, event interface{}) error {
	return Publish(TopicEvents, eventType, event)
}

// Consumer creates a new Kafka consumer
func NewConsumer(topic string, groupID string) *kafka.Reader {
	brokers := strings.Split(config.AppConfig.KafkaBrokers, ",")

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
}

// ConsumeOrders starts consuming order messages
func ConsumeOrders(handler func(key, value []byte) error) {
	groupID := config.AppConfig.KafkaGroupID + "-orders"
	reader := NewConsumer(TopicOrders, groupID)
	defer reader.Close()

	log.Printf("[Kafka] Starting order consumer (group: %s)", groupID)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("[Kafka] Error reading message: %v", err)
			continue
		}

		if err := handler(msg.Key, msg.Value); err != nil {
			log.Printf("[Kafka] Error handling message: %v", err)
		}
	}
}

// ConsumeTrades starts consuming trade messages
func ConsumeTrades(handler func(key, value []byte) error) {
	groupID := config.AppConfig.KafkaGroupID + "-trades"
	reader := NewConsumer(TopicTrades, groupID)
	defer reader.Close()

	log.Printf("[Kafka] Starting trade consumer (group: %s)", groupID)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("[Kafka] Error reading message: %v", err)
			continue
		}

		if err := handler(msg.Key, msg.Value); err != nil {
			log.Printf("[Kafka] Error handling message: %v", err)
		}
	}
}
