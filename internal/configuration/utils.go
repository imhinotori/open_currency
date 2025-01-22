package configuration

import (
	"github.com/charmbracelet/log"
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvBool(key string) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Environment variable %s not found, fallbacking to false", key)
		return false
	}

	rKey, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing environment variable %s", key)
		return false

	}
	return rKey
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
