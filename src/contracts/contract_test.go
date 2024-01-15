package contracts

import (
	"testing"

	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
	"github.com/devdetour/ulysses/test"
	"github.com/stretchr/testify/assert"
)

// Test getting user token associated with contract from models.
func TestGetToken(t *testing.T) {
	test.Setup()
	defer test.TearDown()

	db := database.DB

	// Create user
	user := models.User{
		Username: "test_user",
		Email:    "test_user@test.com",
		Password: "hashed_password_i_promise",
		Names:    "what the heck is this field",
		AuthTokens: []models.ExternalAuthToken{
			{
				UserId: 1,
				Text:   "real-valid-token",
				Scope:  "read",
				Source: models.StravaTokenSource,
			},
		},
	}
	db.Create(&user)
	token, err := GetToken(1, models.StravaTokenSource)

	assert.Nil(t, err)
	assert.Equal(t, "real-valid-token", token.Text)

}
