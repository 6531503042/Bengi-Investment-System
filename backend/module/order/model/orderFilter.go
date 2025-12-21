package model

// OrderFilter is used for filtering orders in queries
type OrderFilter struct {
	Status string
	Side   string
	Symbol string
}
