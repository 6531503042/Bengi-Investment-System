package seeder

import (
	"context"
	"log"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/auth/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/repository"
)

// DefaultRoles defines the default roles to seed
var DefaultRoles = []model.Role{
	{
		Name: model.RoleAdmin,
		Permissions: []string{
			"users:read", "users:write", "users:delete",
			"roles:read", "roles:write",
			"accounts:read", "accounts:write",
			"portfolios:read", "portfolios:write",
			"orders:read", "orders:write",
			"trades:read", "trades:write",
			"instruments:read", "instruments:write",
		},
	},
	{
		Name: model.RoleTrader,
		Permissions: []string{
			"accounts:read", "accounts:write",
			"portfolios:read", "portfolios:write",
			"orders:read", "orders:write",
			"trades:read",
			"instruments:read",
		},
	},
	{
		Name: model.RoleUser,
		Permissions: []string{
			"accounts:read",
			"portfolios:read",
			"orders:read",
			"trades:read",
			"instruments:read",
		},
	},
}

// SeedRoles seeds default roles if they don't exist
func SeedRoles(ctx context.Context) error {
	roleRepo := repository.NewRoleRepository()

	for _, role := range DefaultRoles {
		exists, err := roleRepo.Exists(ctx, role.Name)
		if err != nil {
			return err
		}

		if !exists {
			role.CreatedAt = time.Now()
			role.UpdatedAt = time.Now()
			if err := roleRepo.Create(ctx, &role); err != nil {
				return err
			}
			log.Printf("✅ Created role: %s", role.Name)
		}
	}

	return nil
}

// RunSeeders runs all seeders
func RunSeeders() {
	ctx := context.Background()

	log.Println("🌱 Running seeders...")

	if err := SeedRoles(ctx); err != nil {
		log.Printf("❌ Failed to seed roles: %v", err)
	}

	if err := SeedInstruments(ctx); err != nil {
		log.Printf("❌ Failed to seed instruments: %v", err)
	}

	// Seed test portfolio for test@test.com
	SeedTestPortfolio()

	log.Println("🌱 Seeders completed")
}
