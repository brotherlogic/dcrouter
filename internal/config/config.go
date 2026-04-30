package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config stores the configuration for dcrouter
type Config struct {
	RouterAddress string `json:"router_address"`
	HostAddress   string `json:"host_address"`
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "dcrouter", "config.json"), nil
}

// ReadConfig reads the configuration from the config file
func ReadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return &cfg, nil
}

// WriteConfig writes the configuration to the config file
func WriteConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
