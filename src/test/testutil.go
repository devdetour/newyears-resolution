package test

import (
	"fmt"
	"os"

	"github.com/devdetour/ulysses/config"
	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
)

func Setup() {
	// Initialize test DB

	// Delete old DB
	dbPath := config.Config("DB_FILE")
	err := os.Remove(dbPath)
	if err != nil {
		fmt.Print("Reset DB!")
	} else {
		fmt.Printf("Couldn't delete DB: %v\n", err)
	}

	database.ConnectDB()

	// contract := models.RecurringStravaContract{
	// 	RecurringContract: models.RecurringContract{
	// 		EvaluationSchedule: "* * * * *",
	// 		EvaluationLookback: time.Hour * 24,
	// 	},
	// }
	// tx := db.Create(&contract)
	fmt.Println("Setup done")
}

func TearDown() {
	db := database.DB

	// Delete old DB
	dbPath := config.Config("DB_FILE")
	err := os.Remove(dbPath)
	if err != nil {
		fmt.Println("Cleanup done")
	}
	fmt.Printf("Failed to cleanup DB: %v\n", err)

	// delete all user
	var users []models.User
	db.Find(&users)
	for _, u := range users {
		db.Delete(&u)
	}

	// delete all tokens
	var tokens []models.ExternalAuthToken
	db.Find(&tokens)
	for _, t := range tokens {
		db.Delete(&t)
	}

	// delete all RecurringStravaContracts
	var contracts []models.ExternalAuthToken
	db.Find(&tokens)
	for _, c := range contracts {
		db.Delete(&c)
	}
}
