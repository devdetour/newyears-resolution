package contracts

import (
	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
)

func GetToken(userId uint, tokenSource models.TokenSource) (models.ExternalAuthToken, error) {
	db := database.DB
	var token models.ExternalAuthToken
	err := db.Find(&token, "User_ID = ?", userId).Error
	if err != nil {
		return token, err
	}
	return token, nil
}
