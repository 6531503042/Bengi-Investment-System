package dto

type (
	QuoteResponse struct {
		Symbol        string  `json:"symbol"`
		Price         float64 `json:"price"`
		Open          float64 `json:"open"`
		High          float64 `json:"high"`
		Low           float64 `json:"low"`
		Close         float64 `json:"close"`
		PreviousClose float64 `json:"previousClose"`
		Volume        int64   `json:"volume"`
		Change        float64 `json:"change"`
		ChangePercent float64 `json:"changePercent"`
		Timestamp     string  `json:"timestamp"`
	}

	MultiQuoteResponse struct {
		Quote []QuoteResponse `json:"quote"`
	}

	TwelveDataQuote struct {
		Symbol        string `json:"symbol"`
		Name          string `json:"name"`
		Exchange      string `json:"exchange"`
		Currency      string `json:"currency"`
		Datetime      string `json:"datetime"`
		Open          string `json:"open"`
		High          string `json:"high"`
		Low           string `json:"low"`
		Close         string `json:"close"`
		Volume        string `json:"volume"`
		PreviousClose string `json:"previous_close"`
		Change        string `json:"change"`
		PercentChange string `json:"percent_change"`
	}
)
