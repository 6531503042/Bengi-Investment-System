package seeder

import (
	"context"
	"log"
	"time"

	accountModel "github.com/bricksocoolxd/bengi-investment-system/module/account/model"
	authModel "github.com/bricksocoolxd/bengi-investment-system/module/auth/model"
	portfolioModel "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sample portfolio positions for test user
var testPositions = []struct {
	Symbol   string
	Quantity float64
	AvgCost  float64
}{
	{"NVDA", 15, 120.50},
	{"TSLA", 8, 245.00},
	{"AAPL", 25, 168.00},
	{"AMZN", 12, 142.50},
	{"PLTR", 100, 18.50},
	{"MSFT", 10, 380.00},
	{"GOOGL", 5, 142.00},
	{"META", 8, 320.00},
}

// SeedTestPortfolio creates portfolio and positions for test@test.com
func SeedTestPortfolio() {
	ctx := context.Background()
	db := database.DB

	// Find test user
	usersCollection := db.Collection(authModel.UserCollection)
	var user authModel.User
	err := usersCollection.FindOne(ctx, bson.M{"email": "test@test.com"}).Decode(&user)
	if err != nil {
		log.Println("⚠️ Test user not found, skipping portfolio seed")
		return
	}

	// Find or create demo account for user
	accountsCollection := db.Collection(accountModel.AccountCollection)
	var account accountModel.Account
	err = accountsCollection.FindOne(ctx, bson.M{"userId": user.ID, "type": "demo"}).Decode(&account)
	if err != nil {
		// Create demo account
		account = accountModel.Account{
			ID:             primitive.NewObjectID(),
			UserID:         user.ID,
			Type:           "demo",
			Balance:        50000.00,
			InitialBalance: 10000.00,
			Currency:       "USD",
			Leverage:       1,
			Status:         "active",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		_, err = accountsCollection.InsertOne(ctx, account)
		if err != nil {
			log.Printf("❌ Failed to create demo account: %v", err)
			return
		}
		log.Println("✅ Demo account created for test@test.com")
	}

	// Check if portfolio already exists
	portfoliosCollection := db.Collection(portfolioModel.PortfolioCollection)
	var existingPortfolio portfolioModel.Portfolio
	err = portfoliosCollection.FindOne(ctx, bson.M{"userId": user.ID}).Decode(&existingPortfolio)
	if err == nil {
		log.Println("⚠️ Portfolio already exists for test@test.com, skipping seed")
		return
	}

	// Create portfolio
	portfolio := portfolioModel.Portfolio{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		AccountID: account.ID,
		Name:      "My Portfolio",
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = portfoliosCollection.InsertOne(ctx, portfolio)
	if err != nil {
		log.Printf("❌ Failed to create portfolio: %v", err)
		return
	}
	log.Printf("✅ Portfolio created for test@test.com (ID: %s)", portfolio.ID.Hex())

	// Create positions
	positionsCollection := db.Collection(portfolioModel.PositionCollection)
	for _, pos := range testPositions {
		position := portfolioModel.Position{
			ID:          primitive.NewObjectID(),
			PortfolioID: portfolio.ID,
			Symbol:      pos.Symbol,
			Quantity:    pos.Quantity,
			AvgCost:     pos.AvgCost,
			TotalCost:   pos.Quantity * pos.AvgCost,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_, err = positionsCollection.InsertOne(ctx, position)
		if err != nil {
			log.Printf("❌ Failed to create position for %s: %v", pos.Symbol, err)
			continue
		}
	}
	log.Printf("✅ Created %d positions for test@test.com", len(testPositions))
}
