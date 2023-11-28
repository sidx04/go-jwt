package initialisers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed to load .env file...")
	}
}
