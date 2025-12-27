package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
)

var (
	ErrQuoteNotFound = errors.New("quote not found")
	ErrAPIError      = errors.New("market data API error")
)

// FinnhubQuote represents the quote response from Finnhub API
type FinnhubQuote struct {
	CurrentPrice  float64 `json:"c"`  // Current price
	Change        float64 `json:"d"`  // Change
	PercentChange float64 `json:"dp"` // Percent change
	High          float64 `json:"h"`  // High price of the day
	Low           float64 `json:"l"`  // Low price of the day
	Open          float64 `json:"o"`  // Open price of the day
	PreviousClose float64 `json:"pc"` // Previous close price
	Timestamp     int64   `json:"t"`  // Unix timestamp
}

type MarketDataService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewMarketDataService() *MarketDataService {
	return &MarketDataService{
		apiKey:  config.AppConfig.FinnhubAPIKey, // Use Finnhub instead of Twelve Data
		baseURL: "https://finnhub.io/api/v1",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *MarketDataService) GetQuote(symbol string) (*model.Quote, error) {
	// Finnhub quote endpoint: /quote?symbol=AAPL&token=YOUR_TOKEN
	url := fmt.Sprintf("%s/quote?symbol=%s&token=%s", s.baseURL, symbol, s.apiKey)

	response, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, ErrAPIError
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var fhQuote FinnhubQuote
	if err := json.Unmarshal(body, &fhQuote); err != nil {
		return nil, err
	}

	// Check if quote is valid (Finnhub returns 0 for invalid symbols)
	if fhQuote.CurrentPrice == 0 && fhQuote.PreviousClose == 0 {
		return nil, ErrQuoteNotFound
	}

	return s.convertFinnhubToQuote(symbol, &fhQuote), nil
}

func (s *MarketDataService) GetMultipleQuotes(symbols []string) ([]model.Quote, error) {
	var quotes []model.Quote
	// Finnhub has generous rate limits (60 calls/minute on free tier)
	for _, symbol := range symbols {
		quote, err := s.GetQuote(symbol)
		if err != nil {
			continue // Skip failed quotes
		}
		quotes = append(quotes, *quote)
	}
	return quotes, nil
}

func (s *MarketDataService) convertFinnhubToQuote(symbol string, fh *FinnhubQuote) *model.Quote {
	timestamp := time.Now()
	if fh.Timestamp > 0 {
		timestamp = time.Unix(fh.Timestamp, 0)
	}

	return &model.Quote{
		Symbol:        symbol,
		Price:         fh.CurrentPrice,
		Open:          fh.Open,
		High:          fh.High,
		Low:           fh.Low,
		Close:         fh.CurrentPrice,
		PreviousClose: fh.PreviousClose,
		Volume:        0, // Finnhub /quote doesn't include volume
		Change:        fh.Change,
		ChangePercent: fh.PercentChange,
		Timestamp:     timestamp,
	}
}

// FinnhubCandleResponse represents the candle response from Finnhub API
type FinnhubCandleResponse struct {
	Close     []float64 `json:"c"` // Close prices
	High      []float64 `json:"h"` // High prices
	Low       []float64 `json:"l"` // Low prices
	Open      []float64 `json:"o"` // Open prices
	Status    string    `json:"s"` // Status: "ok" or "no_data"
	Timestamp []int64   `json:"t"` // Unix timestamps
	Volume    []int64   `json:"v"` // Volume data
}

// Candle represents a single candlestick data point
type Candle struct {
	Time   int64   `json:"time"` // Unix timestamp
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// GetCandles fetches historical candlestick data
// First tries Yahoo Finance (free, no auth), then falls back to synthetic data
func (s *MarketDataService) GetCandles(symbol string, resolution string, from, to int64) ([]Candle, error) {
	// Try Yahoo Finance first (free, no auth, ~2000 req/hr)
	candles, err := s.GetYahooFinanceCandles(symbol, resolution, from, to)
	if err == nil && len(candles) > 0 {
		return candles, nil
	}

	// Fallback to synthetic data based on current quote
	return s.GenerateSyntheticCandles(symbol, resolution, from, to)
}

// YahooFinanceResponse represents the response from Yahoo Finance API
type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// GetYahooFinanceCandles fetches real historical data from Yahoo Finance (free, no API key needed)
func (s *MarketDataService) GetYahooFinanceCandles(symbol string, resolution string, from, to int64) ([]Candle, error) {
	// Convert resolution to Yahoo Finance interval
	// 1m, 2m, 5m, 15m, 30m, 60m, 90m, 1h, 1d, 5d, 1wk, 1mo, 3mo
	interval := "1d" // Default daily
	switch resolution {
	case "1":
		interval = "1m"
	case "5":
		interval = "5m"
	case "15":
		interval = "15m"
	case "30":
		interval = "30m"
	case "60":
		interval = "1h"
	case "D":
		interval = "1d"
	case "W":
		interval = "1wk"
	case "M":
		interval = "1mo"
	}

	// Convert crypto symbols for Yahoo Finance (BTC/USD -> BTC-USD)
	yahooSymbol := strings.ReplaceAll(symbol, "/", "-")

	// Yahoo Finance API URL
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&interval=%s",
		yahooSymbol, from, to, interval,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, ErrAPIError
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var yahooResp YahooFinanceResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, err
	}

	// Check for errors
	if yahooResp.Chart.Error != nil {
		return nil, errors.New(yahooResp.Chart.Error.Description)
	}

	if len(yahooResp.Chart.Result) == 0 || len(yahooResp.Chart.Result[0].Timestamp) == 0 {
		return nil, ErrQuoteNotFound
	}

	result := yahooResp.Chart.Result[0]
	quotes := result.Indicators.Quote[0]

	candles := make([]Candle, len(result.Timestamp))
	for i := range result.Timestamp {
		open := 0.0
		high := 0.0
		low := 0.0
		closePrice := 0.0
		volume := int64(0)

		if i < len(quotes.Open) && quotes.Open[i] != 0 {
			open = quotes.Open[i]
		}
		if i < len(quotes.High) && quotes.High[i] != 0 {
			high = quotes.High[i]
		}
		if i < len(quotes.Low) && quotes.Low[i] != 0 {
			low = quotes.Low[i]
		}
		if i < len(quotes.Close) && quotes.Close[i] != 0 {
			closePrice = quotes.Close[i]
		}
		if i < len(quotes.Volume) {
			volume = quotes.Volume[i]
		}

		candles[i] = Candle{
			Time:   result.Timestamp[i],
			Open:   math.Round(open*100) / 100,
			High:   math.Round(high*100) / 100,
			Low:    math.Round(low*100) / 100,
			Close:  math.Round(closePrice*100) / 100,
			Volume: volume,
		}
	}

	return candles, nil
}

// GenerateSyntheticCandles creates realistic chart data based on current quote
// Used as fallback when Finnhub candle API is not available (free tier)
func (s *MarketDataService) GenerateSyntheticCandles(symbol string, resolution string, from, to int64) ([]Candle, error) {
	// Get current quote to base synthetic data on
	quote, err := s.GetQuote(symbol)
	if err != nil {
		return nil, err
	}

	currentPrice := quote.Price
	if currentPrice == 0 {
		currentPrice = 100 // Default fallback
	}

	// Determine interval based on resolution
	var interval int64
	switch resolution {
	case "1":
		interval = 60
	case "5":
		interval = 5 * 60
	case "15":
		interval = 15 * 60
	case "30":
		interval = 30 * 60
	case "60":
		interval = 60 * 60
	case "D":
		interval = 24 * 60 * 60
	case "W":
		interval = 7 * 24 * 60 * 60
	case "M":
		interval = 30 * 24 * 60 * 60
	default:
		interval = 24 * 60 * 60 // Default to daily
	}

	// Generate candles from 'from' to 'to'
	numCandles := int((to - from) / interval)
	if numCandles < 10 {
		numCandles = 30 // Minimum 30 candles
	}
	if numCandles > 365 {
		numCandles = 365 // Max 1 year of daily candles
	}

	candles := make([]Candle, numCandles)

	// Start from a price that trends up to current price
	// Use change percent to determine trend direction
	trendUp := quote.ChangePercent >= 0
	startPrice := currentPrice
	if trendUp {
		startPrice = currentPrice * 0.85 // Started 15% lower
	} else {
		startPrice = currentPrice * 1.15 // Started 15% higher
	}

	// Generate realistic looking chart
	price := startPrice
	volatility := currentPrice * 0.02 // 2% daily volatility

	for i := 0; i < numCandles; i++ {
		timestamp := from + int64(i)*interval

		// Random walk with trend
		change := (rand.Float64() - 0.45) * volatility // Slight upward bias
		if !trendUp {
			change = (rand.Float64() - 0.55) * volatility // Slight downward bias
		}

		open := price
		close := price + change

		// Ensure we converge to current price at the end
		if i == numCandles-1 {
			close = currentPrice
		}

		// Generate high/low with some randomness
		spread := volatility * rand.Float64()
		high := math.Max(open, close) + spread*0.5
		low := math.Min(open, close) - spread*0.5

		// Generate volume (random but decreasing towards current)
		volume := int64(1000000 + rand.Float64()*5000000)

		candles[i] = Candle{
			Time:   timestamp,
			Open:   math.Round(open*100) / 100,
			High:   math.Round(high*100) / 100,
			Low:    math.Round(low*100) / 100,
			Close:  math.Round(close*100) / 100,
			Volume: volume,
		}

		price = close
	}

	return candles, nil
}

// GetDailyCandles fetches daily candles for the last N days
func (s *MarketDataService) GetDailyCandles(symbol string, days int) ([]Candle, error) {
	to := time.Now().Unix()
	from := time.Now().AddDate(0, 0, -days).Unix()
	return s.GetCandles(symbol, "D", from, to)
}

// GetHourlyCandles fetches hourly candles for the last N hours
func (s *MarketDataService) GetHourlyCandles(symbol string, hours int) ([]Candle, error) {
	to := time.Now().Unix()
	from := time.Now().Add(time.Duration(-hours) * time.Hour).Unix()
	return s.GetCandles(symbol, "60", from, to)
}
