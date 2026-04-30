package engine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFetchMappings(t *testing.T) {
	data := Mappings{
		Containers: map[string]ContainerInfo{
			"test": {Port: 1234},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(data)
	}))
	defer server.Close()

	mappings, err := fetchMappings(server.URL)
	if err != nil {
		t.Fatalf("fetchMappings failed: %v", err)
	}

	if mappings.Containers["test"].Port != 1234 {
		t.Errorf("Expected port 1234, got %d", mappings.Containers["test"].Port)
	}
}

func TestGetMappingsWithCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dcrouter-test")
	if err != nil {
		t.Fatalf("Failed to create tmp dir: %v", tmpDir)
	}
	defer os.RemoveAll(tmpDir)

	oldCachePath := cachePath
	cachePath = filepath.Join(tmpDir, "mappings.json")
	defer func() { cachePath = oldCachePath }()

	data := Mappings{
		Containers: map[string]ContainerInfo{
			"cached": {Port: 5678},
		},
	}
	err = saveToCache(&data)
	if err != nil {
		t.Fatalf("saveToCache failed: %v", err)
	}

	// Mock server that shouldn't be called if cache is valid
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Server should not have been called")
	}))
	defer server.Close()

	// Update mappingsURL temporarily
	oldURL := mappingsURL
	mappingsURL = server.URL
	defer func() { mappingsURL = oldURL }()

	mappings, err := GetMappings()
	if err != nil {
		t.Fatalf("GetMappings failed: %v", err)
	}

	if mappings.Containers["cached"].Port != 5678 {
		t.Errorf("Expected port 5678, got %d", mappings.Containers["cached"].Port)
	}
}

func TestGetMappingsExpiredCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dcrouter-test")
	if err != nil {
		t.Fatalf("Failed to create tmp dir: %v", tmpDir)
	}
	defer os.RemoveAll(tmpDir)

	oldCachePath := cachePath
	cachePath = filepath.Join(tmpDir, "mappings.json")
	defer func() { cachePath = oldCachePath }()

	data := Mappings{
		Containers: map[string]ContainerInfo{
			"old": {Port: 1111},
		},
	}
	err = saveToCache(&data)
	if err != nil {
		t.Fatalf("saveToCache failed: %v", err)
	}

	// Backdate the cache file
	err = os.Chtimes(cachePath, time.Now(), time.Now().Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}

	newData := Mappings{
		Containers: map[string]ContainerInfo{
			"new": {Port: 2222},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(newData)
	}))
	defer server.Close()

	// Update mappingsURL temporarily
	oldURL := mappingsURL
	mappingsURL = server.URL
	defer func() { mappingsURL = oldURL }()

	mappings, err := GetMappings()
	if err != nil {
		t.Fatalf("GetMappings failed: %v", err)
	}

	if mappings.Containers["new"].Port != 2222 {
		t.Errorf("Expected port 2222, got %d", mappings.Containers["new"].Port)
	}
}
