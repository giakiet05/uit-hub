package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ChaosConfig struct {
	Enabled    bool
	ErrorRate  float64
	MinDelayMS int
	MaxDelayMS int
}

type Config struct {
	Chaos ChaosConfig
}

func NewConfig() *Config {
	loadEnv()

	return &Config{
		Chaos: ChaosConfig{
			Enabled:    getBool("FAKE_CHAOS_ENABLED", false),
			ErrorRate:  getFloat("FAKE_ERROR_RATE", 0),
			MinDelayMS: getInt("FAKE_MIN_DELAY_MS", 0),
			MaxDelayMS: getInt("FAKE_MAX_DELAY_MS", 0),
		},
	}
}

func loadEnv() {
	_ = godotenv.Load()
	_ = godotenv.Load("apps/fake-uit-server/.env")
}

func getBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
