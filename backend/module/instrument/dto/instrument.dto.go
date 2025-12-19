package dto

type (
	CreateInstrumentRequest struct {
		Symbol      string `json:"symbol" validate:"required,min=1,max=20"`
		Name        string `json:"name" validate:"required,min=1,max=100"`
		Type        string `json:"type" validate:"required,oneof=STOCK ETF CRYPTO FOREX"`
		Exchange    string `json:"exchange" validate:"required"`
		Currency    string `json:"currency" validate:"required"`
		Description string `json:"description"`
		LogoURL     string `json:"logoUrl"`
	}

	UpdateInstrumentRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		LogoURL     string `json:"logoUrl"`
		Status      string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE DELISTED"`
	}

	InstrumentResponse struct {
		ID          string `json:"id"`
		Symbol      string `json:"symbol"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Exchange    string `json:"exchange"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
		LogoURL     string `json:"logoUrl,omitempty"`
		Status      string `json:"status"`
	}

	InstrumentListResponse struct {
		Instruments []InstrumentResponse `json:"instruments"`
		Total       int                  `json:"total"`
		Page        int                  `json:"page"`
		Limit       int                  `json:"limit"`
	}

	SearchQuery struct {
		Query    string `query:"q"`
		Type     string `query:"type"`
		Exchange string `query:"exchange"`
		Page     int    `query:"page"`
		Limit    int    `query:"limit"`
	}
)
