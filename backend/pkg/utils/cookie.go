package utils

import (
	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"github.com/gofiber/fiber/v2"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

func SetAuthCookies(c *fiber.Ctx, accessToken string, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(config.AppConfig.JWTExpireDuration.Seconds()), // e.g., 15 minutes
		Secure:   config.AppConfig.Env == "production",              // HTTPS only in prod
		HTTPOnly: true,                                              // ป้องกัน XSS
		SameSite: "Lax",                                             // ป้องกัน CSRF
	})
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh", // จำกัด path
		MaxAge:   60 * 60 * 24 * 7,       // 7 days
		Secure:   config.AppConfig.Env == "production",
		HTTPOnly: true,
		SameSite: "Strict",
	})
}

func ClearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		MaxAge:   -1,
		HTTPOnly: true,
	})
}
