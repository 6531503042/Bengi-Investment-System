package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Initialze() {
	InitBus()
	InitManager()
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

	app.Get("/ws", UpgradeMiddleware())
	app.Get("/ws", websocket.New(Handler))

	app.Get("/ws/stats", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"clients": Manager.ClientCount(),
			"topics":  Bus.GetTopics(),
		})
	})
}
