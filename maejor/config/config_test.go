package config

import (
	"os"
	"path/filepath"
	"testing"
)

const config = ".maejor.yaml"

func locationOfTestFixture(t *testing.T) string {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return pwd
}

func verifyTestConfigFileExist(t *testing.T) {
	configPath := locationOfTestFixture(t)
	configFile := filepath.Join(configPath, config)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Expected: Test config exist Got: %v", err)
	}
}

func removeTestConfigFile(t *testing.T) {
	configPath := locationOfTestFixture(t)
	configFile := filepath.Join(configPath, config)
	if err := os.RemoveAll(configFile); err != nil {
		t.Fatal(err)
	}
}

func TestConfig(t *testing.T) {
	configPath := locationOfTestFixture(t)
	if err := Create(configPath, configPath); err != nil {
		t.Fatal(err)
	}
	verifyTestConfigFileExist(t)
	removeTestConfigFile(t)
}
