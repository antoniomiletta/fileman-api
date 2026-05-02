package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func mustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return val
}

func getDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("env var %q must be a duration, got %q", key, val)
	}
	return duration
}

func getInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	number, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("env var %q must be a number, got %q", key, val)
	}
	return number
}
