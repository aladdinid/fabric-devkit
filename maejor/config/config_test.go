package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aladdinid/fabric-devkit/maejor/config"
)

func locationOfTestFixture(t *testing.T) string {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to find fixture")
	}

	fixture := filepath.Join(pwd, "fixture")
	if _, err := os.Stat(fixture); os.IsNotExist(err) {
		t.Fatal("Unable to find fixture")
	}
	return fixture
}

func createTestConfigFile(t *testing.T) {

	fixture := locationOfTestFixture(t)
	configFile := filepath.Join(fixture, config.ConfigFilename)
	_, err := os.Create(configFile)
	if err != nil {
		t.Fatal(err)
	}
}

func verifyTestConfigFileExist(t *testing.T) {
	configPath := locationOfTestFixture(t)
	configFile := filepath.Join(configPath, config.ConfigFilename)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Expected: Test config exist Got: %v", err)
	}
}

func removeTestConfigFile(t *testing.T) {
	configPath := locationOfTestFixture(t)
	configFile := filepath.Join(configPath, config.ConfigFilename)
	if err := os.RemoveAll(configFile); err != nil {
		t.Fatal(err)
	}
}

func TestCreateConfig(t *testing.T) {
	configPath := locationOfTestFixture(t)
	if err := config.Create(configPath, configPath); err != nil {
		t.Fatal(err)
	}
	verifyTestConfigFileExist(t)
	removeTestConfigFile(t)
}

func TestSearch(t *testing.T) {

	fixture := locationOfTestFixture(t)

	noConfigPath := filepath.Join(fixture, "..")
	result := config.Search(noConfigPath)
	actual := len(result)
	expected := 0
	if expected != actual {
		removeTestConfigFile(t)
		t.Fatalf("Expected: %d config file Got: %d", expected, actual)
	}

	createTestConfigFile(t)
	result = config.Search(fixture)
	actual = len(result)
	expected = 1
	if expected != actual {
		removeTestConfigFile(t)
		t.Fatalf("Expected: %d config file Got: %d", expected, actual)
	}

	removeTestConfigFile(t)
}
