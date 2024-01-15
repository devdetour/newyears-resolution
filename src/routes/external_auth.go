package routes

import (
	"github.com/gofiber/fiber/v2"
)

// TODO maybe this should redirect to frontend & then POST token to external_auth endpt? ... no I don't think so tbh
func ExternalAuthRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("External auth endpoint")
	})
}
