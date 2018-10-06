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
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aladdinid/fabric-devkit/internal/config"
)

func TestMain(m *testing.M) {

	err := config.Initialize("./testdata", "config")
	if err != nil {
		log.Fatalf("Unable to configure config object: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)

}

func TestProjectPath(t *testing.T) {
	expected := "<path to project>/fabric-devkit"
	actual := config.ProjectPath()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestHyperledgerImages(t *testing.T) {
	expected := 6
	result := config.HyperledgerImages()
	actual := len(result)
	if expected != actual {
		t.Fatalf("Expected: %d images Got: %d images", expected, actual)
	}
}

func TestDomain(t *testing.T) {
	expected := "fabric.network"
	actual := config.Domain()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

type fixtureType struct {
	Name string
	ID   string
}

func TestOrganizations(t *testing.T) {

	result := config.Organizations()
	actual := len(result)
	expected := 2

	if expected != actual {
		t.Fatalf("Expected: %d items Got %d", expected, actual)
	}

}
