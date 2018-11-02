// +build unit

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
	"bytes"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	var b bytes.Buffer
	w := io.Writer(&b)

	testPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: Unable to get pwd. Got: %v", err)
	}

	if err := ConfigTemplate.Execute(w, struct {
		ProjectPath string
	}{
		testPath,
	}); err != nil {
		log.Fatalf("Error: unable to generate config writer. Got: %v", err)
	}

	err = initializeByReader(b.Bytes())
	if err != nil {
		log.Fatalf("Unable to configure config object: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)

}

func TestProjectPath(t *testing.T) {
	expected := "fabric-devkit/maejor/svc"
	actual := ProjectPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestNetworkPath(t *testing.T) {
	expected := "fabric-devkit/maejor/svc"
	actual := NetworkPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestCryptoPath(t *testing.T) {
	expected := "fabric-devkit/maejor/svc/network/crypto-config"
	actual := CryptoPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestChannelArtefactPath(t *testing.T) {
	expected := "fabric-devkit/maejor/svc/network/channel-artefacts"
	actual := ChannelArtefactPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestScriptPath(t *testing.T) {
	expected := "fabric-devkit/maejor/svc/network/scripts"
	actual := ScriptPath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestChaincodePath(t *testing.T) {
	expected := "src/github.com/aladdinid/chaincodes"
	actual := ChaincodePath()
	if !strings.Contains(actual, expected) {
		t.Fatalf("Expected: substring %s in %s", expected, actual)
	}
}

func TestHyperledgerImages(t *testing.T) {
	expected := 6
	result := HyperledgerImages()
	actual := len(result)
	if expected != actual {
		t.Fatalf("Expected: %d images Got: %d images", expected, actual)
	}
}

func TestDomain(t *testing.T) {
	expected := "fabric.network"
	actual := Domain()
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("Expected: string value %s Got: %s", expected, actual)
	}
}

func TestConsortiumByName(t *testing.T) {

	t.Run("ConsortiumNotFound", func(t *testing.T) {
		actual := consortiumByName("1")
		expected := ConsortiumSpec{}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}

	})

	t.Run("ChannelByName", func(t *testing.T) {
		actual := channelByName("ChannelOne")
		expected := ChannelSpec{Name: "ChannelOne", Organizations: []string{"Org1", "Org2"}}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v, Got: %v", expected, actual)
		}

	})

	t.Run("ConsortiumByName", func(t *testing.T) {
		actual := consortiumByName("SampleConsortium")
		channelSpecs := []ChannelSpec{
			ChannelSpec{Name: "ChannelOne", Organizations: []string{"Org1", "Org2"}},
			ChannelSpec{Name: "ChannelTwo", Organizations: []string{"Org2"}},
		}
		expected := ConsortiumSpec{Name: "SampleConsortium", ChannelSpecs: channelSpecs}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v, Got: %v", expected, actual)
		}

	})

}

func TestConsortiumSpecs(t *testing.T) {

	channelSpecs := []ChannelSpec{
		ChannelSpec{Name: "ChannelOne", Organizations: []string{"Org1", "Org2"}},
		ChannelSpec{Name: "ChannelTwo", Organizations: []string{"Org2"}},
	}

	expected := []ConsortiumSpec{
		ConsortiumSpec{Name: "SampleConsortium", ChannelSpecs: channelSpecs},
	}
	actual := ConsortiumSpecs()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %v Got: %v", expected, actual)
	}

}

func TestOrgByName(t *testing.T) {

	t.Run("OrgNotFound", func(t *testing.T) {
		actual := orgByName("O1")
		expected := OrgSpec{}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

	t.Run("VerifyType", func(t *testing.T) {
		org := orgByName("Org1")
		value := reflect.ValueOf(&org).Elem()

		expected := 3
		actual := value.NumField()

		if expected != actual {
			t.Fatalf("Expected: %d fields Got: %d", expected, actual)
		}
	})

	t.Run("OrgFound", func(t *testing.T) {
		actual := orgByName("Org1")
		expected := OrgSpec{
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

	expected := []OrgSpec{
		OrgSpec{
			Name:   "Org1",
			ID:     "Org1MSP",
			Anchor: "peer0",
		},
		OrgSpec{
			Name:   "Org2",
			ID:     "Org2MSP",
			Anchor: "peer0",
		},
	}
	actual := OrganizationSpecs()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %v Got: %v", expected, actual)
	}

}

func TestNewNetworkSpec(t *testing.T) {

	spec := NewNetworkSpec()
	value := reflect.ValueOf(*spec)

	expected := 8
	actual := value.NumField()

	if expected != actual {
		t.Fatalf("Expected number of fields %d but actual %d", expected, actual)
	}

}
