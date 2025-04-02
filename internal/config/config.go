package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     int
	HttpPort string `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"localhost:8081"`
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string

	MinConn int
	MaxConn int
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	portStr := os.Getenv("GRPS_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("GRPS_PORT is not set %s", portStr)
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid GRPS_PORT: %w", err)
	}

	return &Config{
		Host:     os.Getenv("GRPS_HOST"),
		Port:     portInt,
		HttpPort: os.Getenv("HTTP_PORT"),
	}, nil
}

func LoadDatabase() (*DatabaseConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	host := os.Getenv("POSTGRES_HOST")
	portStr := os.Getenv("POSTGRES_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("POSTGRES_PORT is not set")
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	minConnStr := os.Getenv("POSTGRES_MIN_CONN")
	if minConnStr == "" {
		return nil, fmt.Errorf("POSTGRES_MIN_CONN is not set")
	}

	minConnInt, err := strconv.Atoi(minConnStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_MIN_CONN: %w", err)
	}

	MaxConnStr := os.Getenv("POSTGRES_MAX_CONN")
	if MaxConnStr == "" {
		return nil, fmt.Errorf("POSTGRES_MAX_CONN is not set")
	}

	MaxConnInt, err := strconv.Atoi(MaxConnStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_MAX_CONN: %w", err)
	}

	return &DatabaseConfig{
		Host:     host,
		Port:     portInt,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
		MinConn:  minConnInt,
		MaxConn:  MaxConnInt,
	}, nil
}
