package config

import (
	"log"
	"os"
	"time"
)

func GetDuration(key string, fallback time.Duration) time.Duration {
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
