package devkit

import (
	"os"
	"strings"
	"testing"

	"github.com/aladdinid/fabric-devkit/maejor/config"
)

func TestConfig(t *testing.T) {
	if err := config.Initialize(); err != nil {
		t.Fatalf("Fail to initialize config: %v", err)
	}
}

func TestProjectPath(t *testing.T) {

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to retrieve present working directory")
	}

	path := config.ProjectPath()

	if strings.Compare(pwd, path) != 0 {
		t.Fatal("Project path in config file not value")
	}

}

func TestDomain(t *testing.T) {

	result := config.Domain()
	if strings.Compare(result, "fabric.network") != 0 {
		t.Fatal("Domain name not the same as config file")
	}
}

func TestOrgnizations(t *testing.T) {
	result := config.Organizations()
	if len(result) != 2 {
		t.Fatalf("Expect: 2 Got: %d", len(result))
	}
}

func TestContainersHyperledger(t *testing.T) {
	result := config.Hyperledger()
	if len(result) != 6 {
		t.Fatalf("Expect: 6 Got: %d", len(result))
	}
}
