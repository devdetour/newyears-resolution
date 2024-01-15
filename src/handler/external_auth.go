package handler

import (
	"fmt"
	"time"

	"github.com/devdetour/ulysses/auth"
	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

// TODO move these functions maybe?
func getJwtUsername(c *fiber.Ctx) string {
	user := c.Locals("user").(*jwt.Token)
	return user.Claims.(jwt.MapClaims)["username"].(string)
}

func getJwtUserId(c *fiber.Ctx) (uint, error) {
	token := c.Locals("user").(*jwt.Token)
	var userId uint
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId = uint(claims["user_id"].(float64))
		if userId == 0 {
			return 0, fmt.Errorf("Couldn't get User ID from token!")
		}
		fmt.Print(userId)
		return userId, nil
	}

	return 0, fmt.Errorf("Failed to get User ID from token!")
}

// TODO make this look for a specific token source
func GetExternalToken(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var token models.ExternalAuthToken
	db.Find(&token, id)
	if token.UserId == 0 { // TODO proper check for this
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No token found for user with id", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Token found", "data": token})
}

// TODO make this look for a specific token source
func GetExternalTokens(c *fiber.Ctx) error {
	// Get user from jwt
	db := database.DB
	username := getJwtUsername(c)
	var user models.User
	db.Find(&user, "Username = ?", username)

	// TODO
	// - make sure username always passed, not email
	// - probably disable this error, too specific
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Could not find user %s", username), "data": nil})
	}
	var token []models.ExternalAuthToken
	db.Find(&token, "User_Id = ?", user.ID)
	if len(token) == 0 { // TODO proper check for this
		return c.Status(200).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("No tokens found for user with id %d", user.ID), "data": "[]"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Token found", "data": token})
}

func CreateExternalToken(c *fiber.Ctx) error {
	db := database.DB
	token := new(models.ExternalAuthToken)

	// Get user associated with the jwt
	username := getJwtUsername(c)
	log.Infof("Username from jwt: %s", username)

	var user models.User
	db.Find(&user, "Username = ?", username)

	// Validate user exists in DB
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Could not find user %s", username), "data": nil})
	}

	// get token source
	if err := c.BodyParser(token); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// For Strava, exchange token
	// TODO check strava token source. pass strava state from frontend for strava auth
	tokenResponse, err := auth.StravaTokenExchange(token.Text)
	if err != nil || tokenResponse == nil {
		log.Error("Failed to get token from Strava!")
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to exchange code for token", "data": err})
	}
	token.Text = tokenResponse.AccessToken

	token.UserId = user.ID
	token.RefreshToken = tokenResponse.RefreshToken
	token.Expires = time.Unix(int64(tokenResponse.ExpiresAt), 0)
	fmt.Print("Storing token: ", token)

	// First, check DB if there is already a token.
	var existingToken models.ExternalAuthToken
	exists := db.Find(&existingToken, "User_Id = ? AND Source = ?", user.ID, token.Source)
	tx := db.Begin()

	if (exists.RowsAffected) > 0 {
		// delete existing
		tx.Delete(&models.ExternalAuthToken{}, "User_Id = ? AND Source = ?", user.ID, token.Source)
		if tx.Error != nil {
			fmt.Print("Failed to delete! Rolling back")
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to replace old token", "data": tx.Error})
		}

	}

	tx.Create(&token)
	if tx.Error != nil {
		fmt.Print("Failed to create new token! Rolling back")
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create new token", "data": tx.Error})
	}

	tx.Commit()
	if tx.Error != nil {
		fmt.Print("Failed to commit! Rolling back")
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to commit transaction", "data": tx.Error})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created token", "data": token})
}

func UpdateExternalToken(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.DB
	var user models.User

	db.First(&user, id)
	user.Names = uui.Names
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}

func DeleteExternalToken(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})

	}

	if !userExists(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})

	}

	db := database.DB
	var user models.User

	db.First(&user, id)

	db.Delete(&user)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
