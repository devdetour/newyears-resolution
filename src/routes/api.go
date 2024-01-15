package routes

import (
	"github.com/devdetour/ulysses/handler"
	"github.com/devdetour/ulysses/middleware"
	"github.com/gofiber/fiber/v2"
)

func APIRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API endpoint")
	})

	// Contract handlers
	contractsGroup := router.Group("/contracts")

	contractsGroup.Get("/get", middleware.Protected(), handler.GetContractsForUser)

	contractsGroup.Get("/history", middleware.Protected(), handler.GetEvaluationHistory)

	// Create a contract
	contractsGroup.Post("/create", middleware.Protected(), handler.CreateContract)

	contractsGroup.Post("/delete", middleware.Protected(), handler.DeleteContract)

	contractsGroup.Post("/update", func(c *fiber.Ctx) error {
		return c.SendString("Update a contract")
	})

	// Data handlers
	dataGroup := router.Group("/data")
	dataGroup.Get("strava", middleware.Protected(), handler.GetStravaDataForUser)
}
