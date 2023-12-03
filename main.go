package main

import (
	"kerja-praktek/config"
	"log"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	router := config.SetupRouter()
	err := router.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
