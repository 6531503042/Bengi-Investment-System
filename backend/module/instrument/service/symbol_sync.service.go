package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/repository"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
)

// FinnhubSymbol represents a symbol from Finnhub API
type FinnhubSymbol struct {
	Currency    string `json:"currency"`
	Description string `json:"description"`
	DisplayName string `json:"displaySymbol"`
	FIGI        string `json:"figi"`
	ISIN        string `json:"isin"`
	MIC         string `json:"mic"`
	ShareClass  string `json:"shareClassFIGI"`
	Symbol      string `json:"symbol"`
	Symbol2     string `json:"symbol2"`
	Type        string `json:"type"`
}

// CryptoSymbol represents a crypto symbol from Finnhub
type CryptoSymbol struct {
	Description   string `json:"description"`
	DisplaySymbol string `json:"displaySymbol"`
	Symbol        string `json:"symbol"`
}

// SymbolSyncService handles syncing all symbols from Finnhub to database
type SymbolSyncService struct {
	repo    *repository.InstrumentRepository
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewSymbolSyncService(repo *repository.InstrumentRepository) *SymbolSyncService {
	return &SymbolSyncService{
		repo:    repo,
		apiKey:  config.AppConfig.FinnhubAPIKey,
		baseURL: "https://finnhub.io/api/v1",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchAllUSSymbols fetches all US stock symbols from Finnhub
func (s *SymbolSyncService) FetchAllUSSymbols() ([]FinnhubSymbol, error) {
	url := fmt.Sprintf("%s/stock/symbol?exchange=US&token=%s", s.baseURL, s.apiKey)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch US symbols: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("finnhub API error: %s - %s", resp.Status, string(body))
	}

	var symbols []FinnhubSymbol
	if err := json.NewDecoder(resp.Body).Decode(&symbols); err != nil {
		return nil, fmt.Errorf("failed to decode symbols: %w", err)
	}

	return symbols, nil
}

// FetchCryptoSymbols fetches crypto symbols from Finnhub
func (s *SymbolSyncService) FetchCryptoSymbols(exchange string) ([]CryptoSymbol, error) {
	// Finnhub crypto exchanges: BINANCE, COINBASE, KRAKEN, etc.
	url := fmt.Sprintf("%s/crypto/symbol?exchange=%s&token=%s", s.baseURL, exchange, s.apiKey)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch crypto symbols: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("finnhub crypto API error: %s", resp.Status)
	}

	var symbols []CryptoSymbol
	if err := json.NewDecoder(resp.Body).Decode(&symbols); err != nil {
		return nil, fmt.Errorf("failed to decode crypto symbols: %w", err)
	}

	return symbols, nil
}

// GetLogoURL generates a logo URL for a symbol
func GetLogoURL(symbol, instrumentType string) string {
	// Use Clearbit for companies, or crypto-specific logos
	switch instrumentType {
	case "Crypto":
		// Extract base currency (BTC from BTC/USD)
		parts := strings.Split(symbol, "/")
		if len(parts) > 0 {
			base := strings.ToLower(parts[0])
			return fmt.Sprintf("https://cryptologos.cc/logos/%s-logo.png", getCryptoFullName(base))
		}
	default:
		// Use various free logo sources with fallback
		cleanSymbol := strings.ToLower(strings.ReplaceAll(symbol, ".", ""))
		return fmt.Sprintf("https://logo.clearbit.com/%s.com", cleanSymbol)
	}
	return ""
}

// getCryptoFullName maps crypto symbols to full names for logo URLs
func getCryptoFullName(symbol string) string {
	mapping := map[string]string{
		"btc":  "bitcoin-btc",
		"eth":  "ethereum-eth",
		"sol":  "solana-sol",
		"xrp":  "xrp-xrp",
		"doge": "dogecoin-doge",
		"ada":  "cardano-ada",
		"bnb":  "bnb-bnb",
		"avax": "avalanche-avax",
		"dot":  "polkadot-new-dot",
		"link": "chainlink-link",
		"uni":  "uniswap-uni",
		"atom": "cosmos-atom",
		"ltc":  "litecoin-ltc",
		"shib": "shiba-inu-shib",
		"pepe": "pepe-pepe",
	}
	if name, ok := mapping[symbol]; ok {
		return name
	}
	return symbol
}

// SyncAllSymbols syncs all symbols from Finnhub to database
func (s *SymbolSyncService) SyncAllSymbols(ctx context.Context) error {
	log.Println("[SymbolSync] Starting full symbol sync from Finnhub...")

	// 1. Fetch all US symbols
	usSymbols, err := s.FetchAllUSSymbols()
	if err != nil {
		log.Printf("[SymbolSync] Warning: failed to fetch US symbols: %v", err)
	} else {
		log.Printf("[SymbolSync] Fetched %d US symbols from Finnhub", len(usSymbols))
	}

	// 2. Fetch crypto symbols from Coinbase (most accessible)
	cryptoSymbols, err := s.FetchCryptoSymbols("COINBASE")
	if err != nil {
		log.Printf("[SymbolSync] Warning: failed to fetch crypto symbols: %v", err)
	} else {
		log.Printf("[SymbolSync] Fetched %d crypto symbols from Coinbase", len(cryptoSymbols))
	}

	// 3. Convert and filter symbols
	var instruments []*model.Instrument
	addedSymbols := make(map[string]bool)

	// Process US stocks and ETFs
	for _, sym := range usSymbols {
		// Skip if already added or empty
		if sym.Symbol == "" || addedSymbols[sym.Symbol] {
			continue
		}

		// Filter out weird symbols (prefer common stocks and ETFs)
		if !isValidSymbol(sym.Symbol) {
			continue
		}

		instrumentType := "Stock"
		if sym.Type == "ETP" || sym.Type == "ETF" {
			instrumentType = "ETF"
		}

		instruments = append(instruments, &model.Instrument{
			Symbol:      sym.Symbol,
			Name:        cleanDescription(sym.Description),
			Type:        model.InstrumentType(instrumentType),
			LogoURL:     GetLogoURLForSymbol(sym.Symbol, instrumentType),
			Description: sym.Description,
			Status:      model.InstrumentStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		addedSymbols[sym.Symbol] = true
	}

	// Process crypto symbols
	for _, sym := range cryptoSymbols {
		// Only include USD pairs
		if !strings.HasSuffix(sym.DisplaySymbol, "/USD") {
			continue
		}

		symbol := sym.DisplaySymbol
		if addedSymbols[symbol] {
			continue
		}

		instruments = append(instruments, &model.Instrument{
			Symbol:      symbol,
			Name:        cleanCryptoDescription(sym.Description),
			Type:        model.InstrumentTypeCrypto,
			LogoURL:     GetLogoURLForSymbol(symbol, "Crypto"),
			Description: sym.Description,
			Status:      model.InstrumentStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		addedSymbols[symbol] = true
	}

	log.Printf("[SymbolSync] Total instruments to sync: %d", len(instruments))

	// 4. Bulk upsert to database
	if len(instruments) > 0 {
		inserted, updated, err := s.repo.BulkUpsertInstruments(ctx, instruments)
		if err != nil {
			return fmt.Errorf("failed to upsert instruments: %w", err)
		}
		log.Printf("[SymbolSync] Sync complete - Inserted: %d, Updated: %d", inserted, updated)
	}

	return nil
}

// isValidSymbol filters out weird/invalid symbols
func isValidSymbol(symbol string) bool {
	// Skip symbols with special characters (except common ones)
	if strings.ContainsAny(symbol, "^$#@!") {
		return false
	}

	// Skip really long symbols (usually obscure)
	if len(symbol) > 6 {
		return false
	}

	// Skip symbols that are all numbers
	allNumbers := true
	for _, c := range symbol {
		if c < '0' || c > '9' {
			allNumbers = false
			break
		}
	}
	if allNumbers {
		return false
	}

	return true
}

// cleanDescription cleans up company description
func cleanDescription(desc string) string {
	// Remove common suffixes
	desc = strings.TrimSuffix(desc, " - Common Stock")
	desc = strings.TrimSuffix(desc, " Common Stock")
	desc = strings.TrimSuffix(desc, " Inc")
	desc = strings.TrimSuffix(desc, " Inc.")
	desc = strings.TrimSuffix(desc, " Corp")
	desc = strings.TrimSuffix(desc, " Corp.")
	desc = strings.TrimSuffix(desc, " Corporation")
	desc = strings.TrimSuffix(desc, " Ltd")
	desc = strings.TrimSuffix(desc, " Ltd.")
	desc = strings.TrimSuffix(desc, " Limited")
	return strings.TrimSpace(desc)
}

// cleanCryptoDescription cleans crypto name
func cleanCryptoDescription(desc string) string {
	// "Bitcoin" -> "Bitcoin"
	parts := strings.Split(desc, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return desc
}

// StartPeriodicSync starts a background job to sync symbols periodically
func (s *SymbolSyncService) StartPeriodicSync(ctx context.Context, interval time.Duration) {
	go func() {
		// Do initial sync
		if err := s.SyncAllSymbols(ctx); err != nil {
			log.Printf("[SymbolSync] Initial sync failed: %v", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("[SymbolSync] Stopping periodic sync")
				return
			case <-ticker.C:
				if err := s.SyncAllSymbols(ctx); err != nil {
					log.Printf("[SymbolSync] Periodic sync failed: %v", err)
				}
			}
		}
	}()
}
