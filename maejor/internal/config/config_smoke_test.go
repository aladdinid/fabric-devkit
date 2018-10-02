// +build smoke

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aladdinid/fabric-devkit/maejor/internal/config"
)

func fixtureCreateConfigFile(t *testing.T) func() {

	t.Helper()

	fixture := filepath.Join(".", "testdata")
	configFile := filepath.Join(fixture, config.ConfigFilename)
	_, err := os.Create(configFile)
	if err != nil {
		t.Fatal(err)
	}
	return func() { os.Remove(configFile) }
}

func fixtureConfigFileExist(t *testing.T) {

	t.Helper()

	configPath := filepath.Join(".", "testdata")
	configFile := filepath.Join(configPath, config.ConfigFilename)

	defer func() {
		if err := os.Remove(configFile); err != nil {
			t.Fatalf("Expected: to remove test config Got: %v", err)
		}
	}()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Expected: Test config exist Got: %v", err)
	}

}

func TestCreateConfig(t *testing.T) {
	configPath := filepath.Join(".", "testdata")
	if err := config.Create(configPath, configPath); err != nil {
		t.Fatal(err)
	}
	fixtureConfigFileExist(t)
}

func TestSearch(t *testing.T) {

	fixture := filepath.Join(".", "testdata")

	t.Run("NoConfig", func(t *testing.T) {
		noConfigPath := filepath.Join(fixture, "..")
		result := config.Search(noConfigPath)
		actual := len(result)
		expected := 0
		if expected != actual {
			t.Fatalf("Expected: %d config file Got: %d", expected, actual)
		}
	})

	t.Run("ConfigExists", func(t *testing.T) {
		removeConfigFunc := fixtureCreateConfigFile(t)
		defer removeConfigFunc()
		result := config.Search(fixture)
		actual := len(result)
		expected := 1
		if expected != actual {
			t.Fatalf("Expected: %d config file Got: %d", expected, actual)
		}
	})

}
