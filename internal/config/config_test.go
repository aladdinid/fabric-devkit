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
	"reflect"
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

func TestNetworkPath(t *testing.T) {
	expected := "<path to project>/fabric-devkit/network"
	actual := config.NetworkPath()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestCryptoPath(t *testing.T) {
	expected := "<path to project>/fabric-devkit/network/crypto-config"
	actual := config.CryptoPath()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestChannelArtefactPath(t *testing.T) {
	expected := "<path to project>/fabric-devkit/network/channel-artefacts"
	actual := config.ChannelArtefactPath()
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

func TestChannelName(t *testing.T) {
	expected := "TwoOrg"
	actual := config.ChannelName()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestConsortium(t *testing.T) {
	expected := "SampleConsortum"
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
