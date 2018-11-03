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

package svc

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Config services tests
func tfixtureCreateConfigFile(t *testing.T) func() {

	t.Helper()

	fixture := filepath.Join(".")
	configFile := filepath.Join(fixture, configFilename)
	_, err := os.Create(configFile)
	if err != nil {
		t.Fatal(err)
	}
	return func() { os.Remove(configFile) }
}

func tfixtureConfigFileExist(t *testing.T) {

	t.Helper()

	configPath := filepath.Join(".")
	configFile := filepath.Join(configPath, configFilename)

	defer func() {
		if err := os.Remove(configFile); err != nil {
			t.Fatalf("Expected: to remove test config Got: %v", err)
		}
	}()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Expected: Test config exist Got: %v", err)
	}

}

func TestConfig(t *testing.T) {

	t.Run("Create(configPath, projectPath)", func(t *testing.T) {
		configPath := filepath.Join(".")
		if err := Create(configPath, configPath); err != nil {
			t.Fatal(err)
		}
		tfixtureConfigFileExist(t)
	})

	t.Run("SearchConfigFile()", func(t *testing.T) {

		fixture := filepath.Join(".")

		t.Run("NoConfig", func(t *testing.T) {
			noConfigPath := filepath.Join(fixture, "..")
			result := SearchConfigFile(noConfigPath)
			actual := len(result)
			expected := 0
			if expected != actual {
				t.Fatalf("Expected: %d config file Got: %d", expected, actual)
			}
		})

		t.Run("ConfigExists", func(t *testing.T) {
			removeConfigFunc := tfixtureCreateConfigFile(t)
			defer removeConfigFunc()
			result := SearchConfigFile(fixture)
			actual := len(result)
			expected := 1
			if expected != actual {
				t.Fatalf("Expected: %d config file Got: %d", expected, actual)
			}
		})

	})

}

// Docker services tests
func tfixtureRemoveImage(t *testing.T) {
	ids, err := searchImage("hello-world")
	if err != nil {
		t.Fatal("image is not found")
	}

	_, err = RemoveImage(ids[0])
	if err != nil {
		t.Fatalf("Error not expected. Got: %v", err)
	}
}

func tfixturePullDummyImage(t *testing.T) {
	reader, err := pullImage(tfixtureHelloWorldLinux)
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	io.Copy(os.Stdout, reader)
}

func tfixtureTagTargetAsSomething(source string) string {
	result := strings.Split(source, ":")
	if len(result) == 2 {
		result[1] = "something"
	}
	return strings.Join(result, ":")
}

func tfixtureLocation(t *testing.T) string {

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return pwd

}

func tfixturePullFabricTools(t *testing.T) {
	reader, err := pullImage("hyperledger/fabric-tools:x86_64-1.1.0")
	if err != nil {
		t.Fatal("Unable to pull fabric-tools:x86_64-1.1.0")
	}

	io.Copy(os.Stdout, reader)

}

func tfixtureTagFabricTools(t *testing.T) {
	err := tagImage("hyperledger/fabric-tools:x86_64-1.1.0", TargetTagAsLatest)
	if err != nil {
		t.Fatalf("Expected no error. Got %v", err)
	}
}

const tfixtureHelloWorldLinux = "hello-world:linux"
const tfixtureHelloWorldLatest = "hello-world"

func TestDocker(t *testing.T) {

	t.Run("BasicOperations", func(t *testing.T) {

		t.Run("PullImage", func(t *testing.T) {
			reader, err := pullImage(tfixtureHelloWorldLinux)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			io.Copy(os.Stdout, reader)
		})

		t.Run("PullImages", func(t *testing.T) {
			if err := PullImages([]string{tfixtureHelloWorldLatest}); err != nil {
				t.Fatalf("Expected no error. Got: %v", err)
			}
		})

		t.Run("SearchImage", func(t *testing.T) {
			if _, err := searchImage("hello-world:*"); err != nil {
				t.Fatalf("Expected: no error Got: %v", err)
			}

			result, err := searchImage("hello-world:*")
			if err != nil {
				t.Fatalf("Expected: no error Got: %v", err)
			}

			if len(result) != 1 {
				t.Fatalf("Expected: 1 Got %d", len(result))
			}
		})

		t.Run("RemoveImage", func(t *testing.T) {
			ids, err := searchImage("hello-world")
			if err != nil {
				t.Fatal("image is not found")
			}

			deleted, err := RemoveImage(ids[0])
			if err != nil {
				t.Fatalf("Expected: no err Got: %v", err)
			}

			if len(deleted) != 6 {
				t.Fatalf("Expected: 6 Got: %d", len(deleted))
			}
		})

	})

	t.Run("TaggingOperations", func(t *testing.T) {
		t.Run("TagImage", func(t *testing.T) {
			tfixturePullDummyImage(t)

			err := tagImage(tfixtureHelloWorldLinux, TargetTagAsLatest)
			if err != nil {
				t.Fatalf("Error not expected. Got: %v", err)
			}
		})

		t.Run("TagImages", func(t *testing.T) {

			t.Helper()
			sources := []string{tfixtureHelloWorldLinux, tfixtureHelloWorldLatest}
			if err := TagImages(sources, tfixtureTagTargetAsSomething); err != nil {
				t.Fatalf("Expected no err. Got %v", err)
			}

			tfixtureRemoveImage(t)
		})
	})

	t.Run("RunCryptoConfigContainer", func(t *testing.T) {
		tfixturePullFabricTools(t)
		tfixtureTagFabricTools(t)
		err := RunCryptoConfigContainer(tfixtureLocation(t), "container_1", "hyperledger/fabric-tools", []string{"which", "cryptogen"})
		if err != nil {
			t.Fatal(err)
		}

		err = RunCryptoConfigContainer(tfixtureLocation(t), "container_2", "hyperledger/fabric-tools", []string{"which", "configtxgen"})
		if err != nil {
			t.Fatal(err)
		}
	})

}
