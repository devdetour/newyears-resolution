package handler

import (
	"fmt"

	"github.com/devdetour/ulysses/connector"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetStravaDataForUser(c *fiber.Ctx) error {
	// Make sure token valid
	// TODO make all this a helper methodddd
	token := c.Locals("user").(*jwt.Token)

	var userId uint
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId = uint(claims["user_id"].(float64))
		if userId == 0 {
			fmt.Errorf("Must have username")
		}
	}

	fmt.Print(userId) // Find tokens for user
	data, err := connector.GetStravaDataForUserId(userId)
	if err != nil {
		fmt.Errorf("Error getting strava data for user %d: %v", userId, err)
	}

	fmt.Print(data)

	return c.Status(200).JSON(fiber.Map{"status": "ok", "message": "Got data", "data": data})
}
