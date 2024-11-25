package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCConfig
	StorageConfig
	TokenTTL time.Duration
	Secret   string
}

type GRPCConfig struct {
	Port    string
	Timeout time.Duration
}

type StorageConfig struct {
	StorageUser     string
	StoragePass     string
	StorageHost     string
	StoragePort     string
	StorageDatabase string
}

func getEnvWithDefault(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func MustLoad() *Config {
	err := godotenv.Load("./config.env")
	if err != nil {
		panic(err)
	}
	minutes, err := strconv.Atoi(getEnvWithDefault("TOKEN_TTL", "60"))
	if err != nil {
		panic(err)
	}
	cfg := Config{
		Secret:   getEnvWithDefault("SECRET", "secret"),
		TokenTTL: time.Duration(minutes) * time.Minute,
		GRPCConfig: GRPCConfig{
			Port:    getEnvWithDefault("GRPC_PORT", "44044"),
			Timeout: 1 * time.Hour,
		},
		StorageConfig: StorageConfig{
			StorageUser:     getEnvWithDefault("STORAGE_USER", "postgres"),
			StoragePass:     getEnvWithDefault("STORAGE_PASSWORD", "postgres"),
			StorageHost:     getEnvWithDefault("STORAGE_HOST", "localhost"),
			StoragePort:     getEnvWithDefault("STORAGE_PORT", "5432"),
			StorageDatabase: getEnvWithDefault("STORAGE_DB", "postgres"),
		},
	}

	return &cfg
}
