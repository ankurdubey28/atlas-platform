package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Println(".env not found, using environment variables")
		}
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
