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
