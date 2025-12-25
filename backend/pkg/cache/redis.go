package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// Initialize connects to Redis
func Initialize() error {
	opt, err := redis.ParseURL(config.AppConfig.RedisURI)
	if err != nil {
		log.Printf("[Redis] Failed to parse URI: %v", err)
		return err
	}

	client = redis.NewClient(opt)

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[Redis] Failed to connect: %v", err)
		return err
	}

	log.Println("✅ Connected to Redis")
	return nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// IsConnected returns true if Redis is connected
func IsConnected() bool {
	if client == nil {
		return false
	}
	return client.Ping(ctx).Err() == nil
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// ==================== String Operations ====================

// Set stores a string value with expiration
func Set(key string, value string, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a string value
func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Delete removes a key
func Delete(key string) error {
	return client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func Exists(key string) (bool, error) {
	result, err := client.Exists(ctx, key).Result()
	return result > 0, err
}

// ==================== JSON Operations ====================

// SetJSON stores a struct as JSON with expiration
func SetJSON(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, data, expiration).Err()
}

// GetJSON retrieves and unmarshals a JSON value
func GetJSON(key string, dest interface{}) error {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// ==================== TTL Operations ====================

// TTL returns the remaining time to live for a key
func TTL(key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}

// Expire sets a new expiration on a key
func Expire(key string, expiration time.Duration) error {
	return client.Expire(ctx, key, expiration).Err()
}

// ==================== Increment Operations ====================

// Incr increments a counter
func Incr(key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// IncrBy increments by a specific amount
func IncrBy(key string, value int64) (int64, error) {
	return client.IncrBy(ctx, key, value).Result()
}

// ==================== Hash Operations ====================

// HSet sets a hash field
func HSet(key, field string, value interface{}) error {
	return client.HSet(ctx, key, field, value).Err()
}

// HGet gets a hash field
func HGet(key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

// HGetAll gets all hash fields
func HGetAll(key string) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields
func HDel(key string, fields ...string) error {
	return client.HDel(ctx, key, fields...).Err()
}

// ==================== List Operations ====================

// LPush prepends values to a list
func LPush(key string, values ...interface{}) error {
	return client.LPush(ctx, key, values...).Err()
}

// RPush appends values to a list
func RPush(key string, values ...interface{}) error {
	return client.RPush(ctx, key, values...).Err()
}

// LRange returns a range of list elements
func LRange(key string, start, stop int64) ([]string, error) {
	return client.LRange(ctx, key, start, stop).Result()
}

// LLen returns the length of a list
func LLen(key string) (int64, error) {
	return client.LLen(ctx, key).Result()
}

// ==================== Set Operations ====================

// SAdd adds members to a set
func SAdd(key string, members ...interface{}) error {
	return client.SAdd(ctx, key, members...).Err()
}

// SMembers returns all members of a set
func SMembers(key string) ([]string, error) {
	return client.SMembers(ctx, key).Result()
}

// SIsMember checks if a value is a member of a set
func SIsMember(key string, member interface{}) (bool, error) {
	return client.SIsMember(ctx, key, member).Result()
}

// SRem removes members from a set
func SRem(key string, members ...interface{}) error {
	return client.SRem(ctx, key, members...).Err()
}

// ==================== Keys Pattern ====================

// Keys returns all keys matching a pattern (use cautiously in production)
func Keys(pattern string) ([]string, error) {
	return client.Keys(ctx, pattern).Result()
}

// FlushDB removes all keys from the current database (use with caution!)
func FlushDB() error {
	return client.FlushDB(ctx).Err()
}
