// +build smoke

/*
Copyright 2018 Aladdin Blockchain Technologies Ltd
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aladdinid/fabric-devkit/internal/config"
)

func tfixtureCreateConfigFile(t *testing.T) func() {

	t.Helper()

	fixture := filepath.Join(".")
	configFile := filepath.Join(fixture, config.ConfigFilename)
	_, err := os.Create(configFile)
	if err != nil {
		t.Fatal(err)
	}
	return func() { os.Remove(configFile) }
}

func tfixtureConfigFileExist(t *testing.T) {

	t.Helper()

	configPath := filepath.Join(".")
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
	configPath := filepath.Join(".")
	if err := config.Create(configPath, configPath); err != nil {
		t.Fatal(err)
	}
	tfixtureConfigFileExist(t)
}

func TestSearch(t *testing.T) {

	fixture := filepath.Join(".")

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
		removeConfigFunc := tfixtureCreateConfigFile(t)
		defer removeConfigFunc()
		result := config.Search(fixture)
		actual := len(result)
		expected := 1
		if expected != actual {
			t.Fatalf("Expected: %d config file Got: %d", expected, actual)
		}
	})

}
