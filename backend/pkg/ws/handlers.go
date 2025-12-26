package ws

import (
	"log"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/cache"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Initialze() {
	InitBus()
	InitManager()

	// Start Finnhub price stream for real-time prices
	stream := GetPriceStream()
	if err := stream.Start(); err != nil {
		log.Printf("[WS] Price stream error: %v", err)
	}
}

func UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func Handler(c *websocket.Conn) {
	userID := c.Query("userId", "")
	if userID == "" {
		if id := c.Locals("userId"); id != nil {
			userID = id.(string)
		}
	}

	client := NewClient(c, userID)
	Manager.Register(client)

	go client.WritePump()
	client.ReadPump()
}

func RegisterRoutes(app *fiber.App) {
	Initialze()

	// WebSocket endpoints - both /ws and /ws/prices work
	app.Get("/ws", UpgradeMiddleware())
	app.Get("/ws", websocket.New(Handler))
	app.Get("/ws/prices", UpgradeMiddleware())
	app.Get("/ws/prices", websocket.New(Handler))

	app.Get("/ws/stats", func(c *fiber.Ctx) error {
		stream := GetPriceStream()
		return c.JSON(fiber.Map{
			"clients":              Manager.ClientCount(),
			"topics":               Bus.GetTopics(),
			"priceStreamConnected": stream.IsConnected(),
			"subscribedSymbols":    stream.GetSubscribedSymbols(),
		})
	})

	// Cache stats endpoint
	app.Get("/ws/cache/stats", func(c *fiber.Ctx) error {
		return c.JSON(cache.CacheStats())
	})

	app.Post("/ws/subscribe/:symbol", func(c *fiber.Ctx) error {
		symbol := c.Params("symbol")
		GetPriceStream().Subscribe(symbol)
		return c.JSON(fiber.Map{
			"message": "Subscribed to " + symbol,
		})
	})
}
