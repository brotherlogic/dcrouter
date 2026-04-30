package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	repo      = "brotherlogic/dcrouter"
	updateURL = "https://api.github.com/repos/brotherlogic/dcrouter/releases/latest"
	cacheTTL  = 24 * time.Hour
)

type Release struct {
	TagName string `json:"tag_name"`
}

func getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(home, ".cache", "dcrouter", "update.json"), nil
}

func GetLatestRelease() (string, error) {
	path, err := getCachePath()
	if err == nil {
		info, err := os.Stat(path)
		if err == nil {
			if time.Since(info.ModTime()) < cacheTTL {
				data, err := os.ReadFile(path)
				if err == nil {
					var release Release
					if err := json.Unmarshal(data, &release); err == nil {
						return release.TagName, nil
					}
				}
			}
		}
	}

	resp, err := http.Get(updateURL)
	if err != nil {
		return "", fmt.Errorf("could not fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("could not decode latest release: %w", err)
	}

	_ = saveToCache(&release)
	return release.TagName, nil
}

func saveToCache(release *Release) error {
	path, err := getCachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("could not create cache directory: %w", err)
	}

	data, err := json.Marshal(release)
	if err != nil {
		return fmt.Errorf("could not marshal release: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write cache file: %w", err)
	}

	return nil
}

func CheckForUpdate(currentVersion string, force bool) {
	if currentVersion == "dev" && !force {
		return
	}

	var latest string
	var err error
	if force {
		// Clear cache if forced
		path, _ := getCachePath()
		if path != "" {
			os.Remove(path)
		}
	}
	
	latest, err = GetLatestRelease()
	if err != nil {
		if force {
			fmt.Printf("Error checking for update: %v\n", err)
		}
		return
	}

	if latest != currentVersion {
		fmt.Printf("A new version of dcr is available: %s (current: %s)\n", latest, currentVersion)
		fmt.Printf("Run 'go install github.com/brotherlogic/dcrouter/cmd/dcr@latest' to update.\n")
	} else if force {
		fmt.Printf("You are on the latest version (%s).\n", currentVersion)
	}
}
