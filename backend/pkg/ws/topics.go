package ws

import "fmt"

// Topic prefixes
const (
	TopicPricePrefix     = "price:"     // price:AAPL, price:GOOG
	TopicOrderPrefix     = "order:"     // order:userId
	TopicTradePrefix     = "trade:"     // trade:userId
	TopicPortfolioPrefix = "portfolio:" // portfolio:userId
)

// Topic constructors
func TopicPrice(symbol string) string {
	return TopicPricePrefix + symbol
}

func TopicOrder(userID string) string {
	return TopicOrderPrefix + userID
}

func TopicTrade(userID string) string {
	return TopicTradePrefix + userID
}

func TopicPortfolio(userID string) string {
	return TopicPortfolioPrefix + userID
}

// ValidateTopic checks if topic is valid
func ValidateTopic(topic string) bool {
	if len(topic) < 3 {
		return false
	}
	// Only allow known prefixes
	prefixes := []string{
		TopicPricePrefix,
		TopicOrderPrefix,
		TopicTradePrefix,
		TopicPortfolioPrefix,
	}
	for _, prefix := range prefixes {
		if len(topic) > len(prefix) && topic[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// IsUserTopic checks if topic is user-specific (requires auth)
func IsUserTopic(topic string) bool {
	return len(topic) > len(TopicOrderPrefix) &&
		(topic[:len(TopicOrderPrefix)] == TopicOrderPrefix ||
			topic[:len(TopicTradePrefix)] == TopicTradePrefix ||
			topic[:len(TopicPortfolioPrefix)] == TopicPortfolioPrefix)
}

// GetUserFromTopic extracts userID from user topic
func GetUserFromTopic(topic string) string {
	prefixes := []string{TopicOrderPrefix, TopicTradePrefix, TopicPortfolioPrefix}
	for _, prefix := range prefixes {
		if len(topic) > len(prefix) && topic[:len(prefix)] == prefix {
			return topic[len(prefix):]
		}
	}
	return ""
}

func TopicPriceAll() string {
	return fmt.Sprintf("%s*", TopicPricePrefix)
}
