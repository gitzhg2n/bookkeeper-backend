package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "bookkeeper.db"
	}
	DB, err = gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	// Auto-migrate all models
	err = DB.AutoMigrate(
		&User{},
		&Household{},
		&Account{},
		&Budget{},
		&Goal{},
		&Investment{},
		&Transaction{},
		&IncomeSource{},
		&HouseholdMember{},
		&BreakupRequest{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database connected and migrated.")
}