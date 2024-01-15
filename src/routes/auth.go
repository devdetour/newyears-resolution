package routes

import (
	"github.com/devdetour/ulysses/handler"
	"github.com/devdetour/ulysses/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	router.Post("/login", handler.Login)
	router.Post("/register", handler.CreateUser)

	// User
	user := router.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Token
	token := router.Group("/token")
	token.Get("/all", middleware.Protected(), handler.GetExternalTokens)
	token.Post("/create", middleware.Protected(), handler.CreateExternalToken) // path is /auth/token/create
}
