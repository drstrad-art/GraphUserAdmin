package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	TenantID     string `json:"tenantId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s\n\nPlease create a config.json file with your Azure AD credentials.\nSee config.json.example for the required format", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w\n\nEnsure the file contains valid JSON", err)
	}

	// Validate required fields with specific error messages
	var missingFields []string
	if cfg.TenantID == "" {
		missingFields = append(missingFields, "tenantId")
	}
	if cfg.ClientID == "" {
		missingFields = append(missingFields, "clientId")
	}
	if cfg.ClientSecret == "" {
		missingFields = append(missingFields, "clientSecret")
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("config file is missing required fields: %v\n\nSee config.json.example for the required format", missingFields)
	}

	return &cfg, nil
}
