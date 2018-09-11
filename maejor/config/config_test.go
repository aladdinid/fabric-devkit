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

func createTestConfigFile(t *testing.T) {
	_, err := os.Create(ConfigFilename)
	if err != nil {
		t.Fatal(err)
	}
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

func TestSearch(t *testing.T) {

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	noConfigPath := filepath.Join(pwd, "..", "consortium")
	result := Search(noConfigPath)
	if len(result) != 0 {
		t.Fatalf("Expected 0, Got: %d", len(result))
	}

	createTestConfigFile(t)
	result = Search(pwd)
	if len(result) == 0 {
		t.Fatalf("Expected: 1, Got: %d", len(result))
	}
	removeTestConfigFile(t)
}
