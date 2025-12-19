package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
)

var (
	ErrQuoteNotFound = errors.New("quote not found")
	ErrAPIError      = errors.New("market data API error")
)

type MarketDataService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewMarketDataService() *MarketDataService {
	return &MarketDataService{
		apiKey:  config.AppConfig.TwelveDataAPIKey,
		baseURL: "https://api.twelvedata.com",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *MarketDataService) GetQuote(symbol string) (*model.Quote, error) {
	url := fmt.Sprintf("%s/quote?symbol=%s&apikey=%s", s.baseURL, symbol, s.apiKey)

	response, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var apiError map[string]interface{}
	if err := json.Unmarshal(body, &apiError); err == nil {
		if apiError["status"] == "error" {
			return nil, ErrAPIError
		}
	}

	var tdQuote dto.TwelveDataQuote
	if err := json.Unmarshal(body, &tdQuote); err != nil {
		return nil, err
	}
	return s.convertToQuote(&tdQuote)
}

func (s *MarketDataService) GetMultipleQuotes(symbols []string) ([]model.Quote, error) {
	var quotes []model.Quote
	// Twelve Data free tier: max 8 symbols per request
	for _, symbol := range symbols {
		quote, err := s.GetQuote(symbol)
		if err != nil {
			continue // Skip failed quotes
		}
		quotes = append(quotes, *quote)
	}
	return quotes, nil
}

func (s *MarketDataService) convertToQuote(td *dto.TwelveDataQuote) (*model.Quote, error) {
	price, _ := strconv.ParseFloat(td.Close, 64)
	open, _ := strconv.ParseFloat(td.Open, 64)
	high, _ := strconv.ParseFloat(td.High, 64)
	low, _ := strconv.ParseFloat(td.Low, 64)
	closePrice, _ := strconv.ParseFloat(td.Close, 64)
	previousClose, _ := strconv.ParseFloat(td.PreviousClose, 64)
	volume, _ := strconv.ParseInt(td.Volume, 10, 64)
	change, _ := strconv.ParseFloat(td.Change, 64)
	changePercent, _ := strconv.ParseFloat(td.PercentChange, 64)
	timestamp, _ := time.Parse("2006-01-02 15:04:05", td.Datetime)
	return &model.Quote{
		Symbol:        td.Symbol,
		Price:         price,
		Open:          open,
		High:          high,
		Low:           low,
		Close:         closePrice,
		PreviousClose: previousClose,
		Volume:        volume,
		Change:        change,
		ChangePercent: changePercent,
		Timestamp:     timestamp,
	}, nil
}
