package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env before main() runs
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
}

func main() {

}
