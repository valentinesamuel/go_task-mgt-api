package pkg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnvVar(key string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return "", err
	}
	return os.Getenv(key), nil
}
