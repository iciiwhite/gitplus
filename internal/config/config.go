package config

import (
	"os"
)

type Config struct {
	GitHubClientID     string
	GitHubClientSecret string
	OpenAIKey          string
	Port               string
}

func Load() (*Config, error) {
	return &Config{
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		OpenAIKey:          os.Getenv("OPENAI_API_KEY"),
		Port:               getEnvOrDefault("GITPLUS_PORT", "8080"),
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}