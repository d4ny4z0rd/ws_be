package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file", err)
	}
}

func GetString(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}
