package model

import "errors"

// Portfolio module errors
var (
	ErrInsufficientShares = errors.New("insufficient shares in position")
	ErrPositionNotFound   = errors.New("position not found")
	ErrPortfolioNotFound  = errors.New("portfolio not found")
)
