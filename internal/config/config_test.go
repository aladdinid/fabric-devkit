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
	"bytes"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aladdinid/fabric-devkit/internal/config"
)

func TestMain(m *testing.M) {

	var b bytes.Buffer
	w := io.Writer(&b)

	testPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: Unable to get pwd. Got: %v", err)
	}

	if err := config.ConfigTemplate.Execute(w, struct {
		ProjectPath string
	}{
		testPath,
	}); err != nil {
		log.Fatalf("Error: unable to generate config writer. Got: %v", err)
	}

	err = config.InitializeByReader(b.Bytes())
	if err != nil {
		log.Fatalf("Unable to configure config object: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)

}

func TestProjectPath(t *testing.T) {
	expected := "fabric-devkit/internal/config"
	actual := config.ProjectPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestNetworkPath(t *testing.T) {
	expected := "fabric-devkit/internal/config/network"
	actual := config.NetworkPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestCryptoPath(t *testing.T) {
	expected := "fabric-devkit/internal/config/network/crypto-config"
	actual := config.CryptoPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestChannelArtefactPath(t *testing.T) {
	expected := "fabric-devkit/internal/config/network/channel-artefacts"
	actual := config.ChannelArtefactPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestScriptPath(t *testing.T) {
	expected := "fabric-devkit/internal/config/network/scripts"
	actual := config.ScriptPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestChaincodePath(t *testing.T) {
	expected := "src/github.com/aladdinid/chaincodes"
	actual := config.ChaincodePath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
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

func TestChannelName(t *testing.T) {
	expected := "TwoOrg"
	actual := config.ChannelName()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestConsortium(t *testing.T) {
	expected := "SampleConsortium"
	actual := config.Consortium()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestOrgByName(t *testing.T) {

	t.Run("OrgNotFound", func(t *testing.T) {
		actual := config.OrgByName("O1")
		expected := config.OrgSpec{}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

	t.Run("VerifyType", func(t *testing.T) {
		org := config.OrgByName("Org1")
		value := reflect.ValueOf(&org).Elem()

		expected := 3
		actual := value.NumField()

		if expected != actual {
			t.Fatalf("Expected: %d fields Got: %d", expected, actual)
		}
	})

	t.Run("OrgFound", func(t *testing.T) {
		actual := config.OrgByName("Org1")
		expected := config.OrgSpec{
			Name:   "Org1",
			ID:     "Org1MSP",
			Anchor: "peer0",
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

}

func TestOrganizationSpecs(t *testing.T) {

	t.Run("FailingNotEqual", func(t *testing.T) {
		expected := []config.OrgSpec{}
		actual := config.OrganizationSpecs()

		if reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

	t.Run("FoundEqual", func(t *testing.T) {
		expected := []config.OrgSpec{
			config.OrgSpec{
				Name:   "Org1",
				ID:     "Org1MSP",
				Anchor: "peer0",
			},
			config.OrgSpec{
				Name:   "Org2",
				ID:     "Org2MSP",
				Anchor: "peer0",
			},
		}
		actual := config.OrganizationSpecs()

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}

	})
}

func TestNewNetworkSpec(t *testing.T) {

	spec := config.NewNetworkSpec()
	value := reflect.ValueOf(*spec)

	expected := value.NumField()
	actual := 9

	if expected != actual {
		t.Fatalf("Expected number of fields %d but actual %d", expected, actual)
	}

}
