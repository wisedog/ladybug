package interfacer

import (
	"testing"
)

// TestLoadConfig validates loading config feature
func TestLoadConfig(t *testing.T) {
	cf := LoadConfig()

	if cf == nil {
		t.Error("config load failed")
	}
}

// TestMode validates GetMode() method
func TestMode(t *testing.T) {
	cf := LoadConfig()

	if cf == nil {
		t.Error("config load failed")
	}

	if cf.GetMode() != "dev" {
		t.Error("Config value is not property")
	}
}

// TestBindAddress validates BindAddress() method
func TestBindAddress(t *testing.T) {
	cf := LoadConfig()

	if cf == nil {
		t.Error("config load failed")
	}

	if cf.GetBindAddress() != "localhost:8000" {
		t.Error("Config value is not property")
	}
}
