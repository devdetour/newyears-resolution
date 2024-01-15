package database

import (
	"github.com/devdetour/ulysses/models"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

func GetExternalAuthTokenForUser(userId uint) (*models.ExternalAuthToken, error) {
	var token models.ExternalAuthToken

	tx := DB.Find(&token, "User_Id = ?", userId) // TODO error check
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &token, nil
}
