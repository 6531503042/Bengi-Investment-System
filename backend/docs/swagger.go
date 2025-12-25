// Package docs Bengi Investment System API
//
// API documentation for Bengi Investment System - a real-time stock trading platform.
//
//	Schemes: http, https
//	BasePath: /api
//	Version: 1.0.0
//	Host: localhost:3000
//
//	SecurityDefinitions:
//	  Bearer:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//	    description: "JWT Bearer token. Format: Bearer {token}"
//
// swagger:meta
package docs

// @title Bengi Investment System API
// @version 1.0
// @description A real-time stock trading platform API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bengi.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Authorization header using Bearer scheme. Example: "Bearer {token}"

// swagger:route POST /auth/login auth loginUser
// Authenticate user and return JWT token.
// responses:
//   200: loginResponse

// swagger:route POST /auth/register auth registerUser
// Register a new user account.
// responses:
//   201: userResponse

// swagger:route GET /instruments instruments listInstruments
// List all available trading instruments.
// responses:
//   200: instrumentsResponse

// swagger:route GET /instruments/{symbol}/quote instruments getQuote
// Get real-time quote for a symbol.
// responses:
//   200: quoteResponse

// swagger:route GET /portfolios portfolios listPortfolios
// List user's portfolios.
// security:
//   - Bearer: []
// responses:
//   200: portfoliosResponse

// swagger:route POST /orders orders createOrder
// Create a new order.
// security:
//   - Bearer: []
// responses:
//   201: orderResponse

// swagger:route GET /trades trades listTrades
// List user's trades.
// security:
//   - Bearer: []
// responses:
//   200: tradesResponse
