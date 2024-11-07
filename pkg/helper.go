package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnvVar(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return "", err
	}
	return os.Getenv(key), nil
}

func GenerateToken(length int) (string, error) {
	// Calculate the number of bytes needed for base64 encoding
	byteLength := (length * 3) / 4
	randomBytes := make([]byte, byteLength)

	// Read random bytes from crypto/rand
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode bytes to base64 string and truncate to the specified length
	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	return randomString[:length], nil
}
