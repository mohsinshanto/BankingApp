package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
}
