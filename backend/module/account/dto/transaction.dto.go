package dto

type (
	TransactionResponse struct {
		ID            string  `json:"id"`
		AccountID     string  `json:"accountId"`
		Type          string  `json:"type"`
		Amount        float64 `json:"amount"`
		BalanceBefore float64 `json:"balanceBefore"`
		BalanceAfter  float64 `json:"balanceAfter"`
		ReferenceType string  `json:"referenceType,omitempty"`
		ReferenceID   *string `json:"referenceId,omitempty"`
		Status        string  `json:"status"`
		Description   string  `json:"description"`
		CreatedAt     string  `json:"createdAt"`
	}

	TransactionFilter struct {
		Type   string `query:"type"`
		Status string `query:"status"`
		Limit  int    `query:"limit"`
		Offset int    `query:"offset"`
	}
)
