//go:build ignore

package main

import (
	"api_orion/db"
	"api_orion/model"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	databaseURL := os.Getenv("DATABASE_URL")
	database := db.Postgres{}
	conn, err := database.ConnectUrl(databaseURL)
	if err != nil {
		panic(err)
	}

	isActive := true
	startDate := time.Now()
	endDate := time.Now().Add(30 * 24 * time.Hour)

	batch := model.Batch{
		Name:      "Test Batch 2026",
		IsActive:  &isActive,
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	if err := conn.Create(&batch).Error; err != nil {
		fmt.Printf("Error creating batch: %v\n", err)
		return
	}
	fmt.Println("Batch created with ID:", batch.ID)
}
