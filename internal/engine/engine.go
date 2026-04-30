package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	mappingsURL = "https://raw.githubusercontent.com/brotherlogic/devcontainer-manager/main/mappings.json"
	cacheTTL    = 5 * time.Minute
	cachePath   = ""
)

// ContainerInfo stores information about a container
type ContainerInfo struct {
	Port int `json:"port"`
}

// Mappings stores the mapping of container names to their info
type Mappings struct {
	Containers map[string]ContainerInfo `json:"containers"`
}

func getCachePath() (string, error) {
	if cachePath != "" {
		return cachePath, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(home, ".cache", "dcrouter", "mappings.json"), nil
}

func fetchMappings(url string) (*Mappings, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch mappings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch mappings: %s", resp.Status)
	}

	var mappings Mappings
	if err := json.NewDecoder(resp.Body).Decode(&mappings); err != nil {
		return nil, fmt.Errorf("could not decode mappings: %w", err)
	}

	return &mappings, nil
}

func saveToCache(mappings *Mappings) error {
	path, err := getCachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("could not create cache directory: %w", err)
	}

	data, err := json.Marshal(mappings)
	if err != nil {
		return fmt.Errorf("could not marshal mappings: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write cache file: %w", err)
	}

	return nil
}

// GetMappings retrieves the mappings from cache or fetches them if expired
func GetMappings() (*Mappings, error) {
	path, err := getCachePath()
	if err == nil {
		info, err := os.Stat(path)
		if err == nil {
			if time.Since(info.ModTime()) < cacheTTL {
				data, err := os.ReadFile(path)
				if err == nil {
					var mappings Mappings
					if err := json.Unmarshal(data, &mappings); err == nil {
						return &mappings, nil
					}
				}
			}
		}
	}

	mappings, err := fetchMappings(mappingsURL)
	if err != nil {
		return nil, err
	}

	_ = saveToCache(mappings) // Best effort
	return mappings, nil
}

// ResolvePort resolves a container name to a port
func ResolvePort(name string) (int, error) {
	mappings, err := GetMappings()
	if err != nil {
		return 0, err
	}

	info, ok := mappings.Containers[name]
	if !ok {
		return 0, fmt.Errorf("container %q not found in mappings", name)
	}

	return info.Port, nil
}
