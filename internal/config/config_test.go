package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	// Create a temporary home directory for testing
	tmpDir, err := os.MkdirTemp("", "dcrouter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override HOME environment variable
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	cfg := &Config{
		RouterAddress: "router.test",
		HostAddress:   "host.test",
	}

	if err := WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	readCfg, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	if readCfg.RouterAddress != cfg.RouterAddress {
		t.Errorf("Expected RouterAddress %s, got %s", cfg.RouterAddress, readCfg.RouterAddress)
	}

	if readCfg.HostAddress != cfg.HostAddress {
		t.Errorf("Expected HostAddress %s, got %s", cfg.HostAddress, readCfg.HostAddress)
	}
}
