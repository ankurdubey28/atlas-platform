package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("C:\\Users\\ASUS\\GolangProjects\\atlas-platform\\.env")
	if err != nil {
		log.Println("No .env file found")
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
