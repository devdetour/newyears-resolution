package server

import (
	"fmt"

	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/routes"
	"github.com/devdetour/ulysses/scheduler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func Run() {
	database.ConnectDB()
	database.CreateSession()
	scheduler.SetupNameToFunction()
	scheduler.SetupScheduler()

	app := fiber.New()

	app.Static("/", "./frontend/build")

	app.Get("/session", func(c *fiber.Ctx) error {
		sess, err := database.Store.Get(c)
		if err != nil {
			panic(err)
		}

		jwt := sess.Get("jwt")
		return c.SendString(fmt.Sprintf("Current JWT: %s", jwt))
	})

	app.Post("/session", func(c *fiber.Ctx) error {
		sess, err := database.Store.Get(c)
		if err != nil {
			panic(err)
		}

		jwt := sess.Get("jwt")
		return c.SendString(fmt.Sprintf("Current JWT: %s", jwt))
	})

	app.Get("/query_activities", func(c *fiber.Ctx) error {
		// token, err := dao.GetToken(DEFAULT_USER)
		// if err != nil {
		// 	return c.SendString("Err")
		// }
		return c.SendString("OK")
	})

	client_id := 11111 // replace with proper one

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/strava_auth": "http://www.strava.com/oauth/authorize?client_id=" + string(client_id) + "&response_type=code&redirect_uri=http://localhost:3000/exchange_token&approval_prompt=force&scope=read",
		},
		StatusCode: 301,
	}))

	// API routes
	routes.APIRoutes(app.Group("/api"))

	// Auth routes for internal
	routes.AuthRoutes(app.Group(("/auth")))

	// Auth routes for external data sources (e.g. strava)
	routes.ExternalAuthRoutes(app.Group("/external_auth"))

	// Contracts routes
	// routes.ContractsRoutes(app.Group("/contracts"))

	app.Static("*", "./frontend/build")

	app.Listen("0.0.0.0:3000")
}
