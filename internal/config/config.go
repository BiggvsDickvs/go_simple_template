package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Jwt_key string
}

func LoadConfig() *Config {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get dir:", err)
	}
	envPath := filepath.Join(basePath, ".env")
	result := godotenv.Load(envPath)
	if result != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Printf("JWT TOKEN: %s", os.Getenv("JWT_KEY"))
	return &Config{Jwt_key: os.Getenv("JWT_KEY")}
}
