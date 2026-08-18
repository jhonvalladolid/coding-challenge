package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv               string
	Port                 string
	MaxMatrixDim         int
	StatisticsAPIURL     string
	StatisticsAPITimeout time.Duration
}

func Load() Config {
	loadDotEnv()

	return Config{
		AppEnv:               getenv("APP_ENV", "development"),
		Port:                 getenv("PORT", "8080"),
		MaxMatrixDim:         getenvInt("MAX_MATRIX_DIM", 200),
		StatisticsAPIURL:     getenv("STATISTICS_API_URL", "http://localhost:3000"),
		StatisticsAPITimeout: getenvDuration("STATISTICS_API_TIMEOUT", 5*time.Second),
	}
}

func (c Config) Validate() error {
	parsed, err := url.Parse(c.StatisticsAPIURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("STATISTICS_API_URL must be an absolute URL")
	}
	if c.StatisticsAPITimeout <= 0 {
		return fmt.Errorf("STATISTICS_API_TIMEOUT must be a positive duration")
	}
	return nil
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func getenvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func loadDotEnv() {
	for _, path := range []string{".env", filepath.Join("go-api", ".env")} {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if key == "" {
				continue
			}
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
			_ = os.Setenv(key, value)
		}

		_ = file.Close()
		return
	}
}
