package config

import (
	"fmt"
	"os"
)

var cfg *Config

func init() {
	cfg = newConfig()
}

func Instance() *Config {
	return cfg
}

func newConfig() *Config {
	return &Config{
		App: AppConfig{
			Port:         require("APP_PORT"),
			AllowOrigins: require("APP_ALLOW_ORIGINS"),
		},
		DB: DBConfig{
			Host:     require("DB_HOST"),
			Port:     require("DB_PORT"),
			Name:     require("DB_NAME"),
			User:     require("DB_USER"),
			Password: require("DB_PASSWORD"),
		},
		Cache: CacheConfig{
			Host:     require("CACHE_HOST"),
			Port:     require("CACHE_PORT"),
			Password: os.Getenv("CACHE_PASSWORD"),
		},
	}
}

func require(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required env var %s is not set", key))
	}

	return val
}
