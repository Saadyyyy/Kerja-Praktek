package main

import (
	"kerja-praktek/config"
	"log"
)

func main() {
	router := config.SetupRouter()
	err := router.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
