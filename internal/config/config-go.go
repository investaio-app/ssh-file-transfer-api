package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	ServerPort        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxRequestSize    int64
	RateLimitRequests int
	RateLimitDuration time.Duration

	// SSH related configuration
	SSHKeyPath  string
	SSHUsername string
	SSHPassword string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Set defaults
	cfg := &Config{
		ServerPort:        "8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxRequestSize:    50 * 1024 * 1024, // 50MB
		RateLimitRequests: 100,
		RateLimitDuration: time.Minute,
		SSHKeyPath:        os.Getenv("SSH_KEY_PATH"),
		SSHUsername:       os.Getenv("SSH_USERNAME"),
		SSHPassword:       os.Getenv("SSH_PASSWORD"),
	}

	// Override defaults with environment variables if provided
	if port := os.Getenv("PORT"); port != "" {
		cfg.ServerPort = port
	}

	if readTimeout := os.Getenv("READ_TIMEOUT"); readTimeout != "" {
		if duration, err := strconv.Atoi(readTimeout); err == nil {
			cfg.ReadTimeout = time.Duration(duration) * time.Second
		}
	}

	if writeTimeout := os.Getenv("WRITE_TIMEOUT"); writeTimeout != "" {
		if duration, err := strconv.Atoi(writeTimeout); err == nil {
			cfg.WriteTimeout = time.Duration(duration) * time.Second
		}
	}

	if maxRequestSize := os.Getenv("MAX_REQUEST_SIZE"); maxRequestSize != "" {
		if size, err := strconv.ParseInt(maxRequestSize, 10, 64); err == nil {
			cfg.MaxRequestSize = size
		}
	}

	if rateLimitRequests := os.Getenv("RATE_LIMIT_REQUESTS"); rateLimitRequests != "" {
		if requests, err := strconv.Atoi(rateLimitRequests); err == nil {
			cfg.RateLimitRequests = requests
		}
	}

	if rateLimitDuration := os.Getenv("RATE_LIMIT_DURATION"); rateLimitDuration != "" {
		if duration, err := strconv.Atoi(rateLimitDuration); err == nil {
			cfg.RateLimitDuration = time.Duration(duration) * time.Second
		}
	}

	// Validate required configurations
	if cfg.SSHKeyPath == "" && cfg.SSHPassword == "" {
		return nil, errors.New("either SSH_KEY_PATH or SSH_PASSWORD must be provided")
	}

	if cfg.SSHUsername == "" {
		return nil, errors.New("SSH_USERNAME is required")
	}

	return cfg, nil
}
