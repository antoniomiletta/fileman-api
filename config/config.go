package config

import "time"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type StorageConfig struct {
	Backend  string
	S3Bucket string
}

type AuthConfig struct {
	JWTSecret   string
	TokenExpiry time.Duration
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			URL:             mustGetEnv("DB_URL"),
			MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 20),
			MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		},
		Storage: StorageConfig{
			Backend:  getEnv("STORAGE_BACKEND", "s3"),
			S3Bucket: getEnv("S3_BUCKET", "./data"),
		},
		Auth: AuthConfig{
			JWTSecret:   mustGetEnv("JWT_SECRET"),
			TokenExpiry: getDuration("TOKEN_EXPIRY", 24*time.Hour),
		},
	}
}
