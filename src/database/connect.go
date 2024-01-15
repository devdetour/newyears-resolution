package database

import (
	"fmt"

	"github.com/devdetour/ulysses/config"
	"github.com/devdetour/ulysses/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB connect to db
func ConnectDB() {
	var err error
	dbFile := config.Config("DB_FILE") // TODO better logic for this. shouldn't need to provide the FULL absolute path.
	// can figure out in parsing logic

	if len(dbFile) == 0 {
		panic("failed to parse database name from config")
	}

	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // TODO maybe only when debug
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Printf("Connection Opened to Database %s\n", dbFile)
	DB.AutoMigrate(&models.User{}, &models.ExternalAuthToken{}, &models.RecurringStravaContract{}, &models.ContractEvaluation{})
	// schema, _ := DB.Migrator().GetTables()

	a, err := DB.Migrator().ColumnTypes(&models.ContractEvaluation{})
	fmt.Println("Table Schema:", a)
	fmt.Println("Database Migrated")
}
