package handler

import (
	"errors"
	"net/mail"
	"time"

	"github.com/devdetour/ulysses/config"
	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Email: e}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Username: u}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"` // username or email
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)
	var userData UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	identity := input.Identity
	pass := input.Password
	usermodels, err := new(models.User), *new(error)

	log.Infof("Trying to log in %s", identity)

	if isEmail(identity) {
		usermodels, err = getUserByEmail(identity)
	} else {
		usermodels, err = getUserByUsername(identity)
	}

	log.Info(usermodels)

	if usermodels == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found", "data": err})
	} else {
		userData = UserData{
			ID:       usermodels.ID,
			Username: usermodels.Username,
			Email:    usermodels.Email,
			Password: usermodels.Password,
		}
	}

	log.Info(userData)

	if !CheckPasswordHash(pass, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userData.Username
	claims["user_id"] = userData.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Store JWT in session storage
	sess, err := database.Store.Get(c)
	if err != nil { // TODO don't panic probably, handle gracefully
		panic(err)
	}
	sess.Set("jwt", t)

	if err := sess.Save(); err != nil {
		panic(err)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
