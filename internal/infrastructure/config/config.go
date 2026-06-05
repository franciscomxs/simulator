package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	Port string
}

// Load reads configuration from environment variables with defaults.
func Load() Config {
	port := os.Getenv("PORT")
	if !isValidPort(port) {
		port = "8080"
	}
	return Config{Port: port}
}

func isValidPort(s string) bool {
	n, err := strconv.Atoi(s)
	return err == nil && n >= 1 && n <= 65535
}
