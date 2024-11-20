package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	RE_TOKEN_AGE  time.Duration // Refresh token max age
	ACC_TOKEN_AGE time.Duration // Access token max age
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Parse refresh token age
	rtageStr := os.Getenv("REFRESH_TOKEN_AGE")
	rtage, err := strconv.Atoi(rtageStr)
	if err != nil {
		log.Printf("Invalid REFRESH_TOKEN_AGE. Using default of 0 minutes. Error: %v", err)
		rtage = 0 // Fallback value
	}
	RE_TOKEN_AGE = time.Duration(rtage) * time.Minute

	// Parse access token age
	atageStr := os.Getenv("ACCESS_TOKEN_AGE")
	atage, err := strconv.Atoi(atageStr)
	if err != nil {
		log.Printf("Invalid ACCESS_TOKEN_AGE. Using default of 0 days. Error: %v", err)
		atage = 0 // Fallback value
	}
	ACC_TOKEN_AGE = time.Duration(atage) * 24 * time.Hour
}
