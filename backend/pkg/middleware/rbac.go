package middleware

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/repository"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/gofiber/fiber/v2"
)

// RoleRequired middleware checks if user has one of the required roles
// Usage: RoleRequired("ADMIN", "TRADER")
func RoleRequired(allowedRoles ...string) fiber.Handler {
	roleRepo := repository.NewRoleRepository()

	return func(c *fiber.Ctx) error {
		// Get roleId from locals (set by AuthRequired middleware)
		roleID, ok := c.Locals("roleId").(string)
		if !ok || roleID == "" {
			return common.Unauthorized(c, "Role not found in token")
		}

		// Get role from database
		role, err := roleRepo.FindByID(c.Context(), roleID)
		if err != nil {
			return common.Unauthorized(c, "Invalid role")
		}

		// Check if user's role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role.Name == allowedRole {
				// Store role name in locals for later use
				c.Locals("roleName", role.Name)
				return c.Next()
			}
		}

		return common.Error(c, fiber.StatusForbidden, "Access denied. Required role: "+joinRoles(allowedRoles))
	}
}

// joinRoles joins role names for error message
func joinRoles(roles []string) string {
	result := ""
	for i, role := range roles {
		if i > 0 {
			result += " or "
		}
		result += role
	}
	return result
}

// GetRoleName gets the role name from context
func GetRoleName(c *fiber.Ctx) string {
	if name, ok := c.Locals("roleName").(string); ok {
		return name
	}
	return ""
}
