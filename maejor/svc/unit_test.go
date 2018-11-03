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

	"github.com/spf13/viper"
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

// Config services tests
func TestConfig(t *testing.T) {

	t.Run("ProjectPath()", func(t *testing.T) {
		expected := "fabric-devkit/maejor/svc"
		actual := ProjectPath()
		if !strings.Contains(actual, expected) {
			t.Fatalf("Expected: substring %s in %s", expected, actual)
		}
	})

	t.Run("NetworkPath()", func(t *testing.T) {
		expected := "fabric-devkit/maejor/svc"
		actual := NetworkPath()
		if !strings.Contains(actual, expected) {
			t.Fatalf("Expected: substring %s in %s", expected, actual)
		}
	})

	t.Run("CryptoPath()", func(t *testing.T) {
		expected := "fabric-devkit/maejor/svc/network/crypto-config"
		actual := CryptoPath()
		if !strings.Contains(actual, expected) {
			t.Fatalf("Expected: substring %s in %s", expected, actual)
		}
	})

	t.Run("ChannelPath()", func(t *testing.T) {
		expected := "fabric-devkit/maejor/svc/network/channel-artefacts"
		actual := ChannelArtefactPath()
		if !strings.Contains(actual, expected) {
			t.Fatalf("Expected: substring %s in %s", expected, actual)
		}
	})

	t.Run("ScriptPath()", func(t *testing.T) {
		expected := "fabric-devkit/maejor/svc/network/scripts"
		actual := ScriptPath()
		if !strings.Contains(actual, expected) {
			t.Fatalf("Expected: substring %s in %s", expected, actual)
		}
	})

	t.Run("ChaincodePath()", func(t *testing.T) {
		expected := "src/github.com/aladdinid/chaincodes"
		actual := ChaincodePath()
		if !strings.Contains(actual, expected) {
			t.Fatalf("Expected: substring %s in %s", expected, actual)
		}
	})

	t.Run("HyperledgerImages()", func(t *testing.T) {
		expected := 6
		result := HyperledgerImages()
		actual := len(result)
		if expected != actual {
			t.Fatalf("Expected: %d images Got: %d images", expected, actual)
		}
	})

	t.Run("Domain()", func(t *testing.T) {
		expected := "fabric.network"
		actual := Domain()
		if strings.Compare(expected, actual) != 0 {
			t.Fatalf("Expected: string value %s Got: %s", expected, actual)
		}
	})

	t.Run("consortiumByName():NotFound", func(t *testing.T) {
		actual := consortiumByName("1")
		expected := ConsortiumSpec{}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

	t.Run("channelByName()", func(t *testing.T) {
		actual := channelByName("ChannelOne")
		expected := ChannelSpec{Name: "ChannelOne", Organizations: []string{"Org1", "Org2"}}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v, Got: %v", expected, actual)
		}
	})

	t.Run("consortiumByName()", func(t *testing.T) {
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

	t.Run("ConsortiumSpecs()", func(t *testing.T) {
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
	})

	t.Run("orgByName():NotFound", func(t *testing.T) {
		actual := orgByName("O1")
		expected := OrgSpec{}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

	t.Run(`orgByName("Org1"):1`, func(t *testing.T) {
		org := orgByName("Org1")
		value := reflect.ValueOf(&org).Elem()

		expected := 3
		actual := value.NumField()

		if expected != actual {
			t.Fatalf("Expected: %d fields Got: %d", expected, actual)
		}
	})

	t.Run(`orgByName("Org1"):2`, func(t *testing.T) {
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

	t.Run("OrganizationSpecs", func(t *testing.T) {
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
	})

}

// Data type services tests
func TestDataTypes(t *testing.T) {

	t.Run("NewNetworkSpec()", func(t *testing.T) {
		spec := NewNetworkSpec()
		value := reflect.ValueOf(*spec)

		expected := 9
		actual := value.NumField()

		if expected != actual {
			t.Fatalf("Expected number of fields %d but actual %d", expected, actual)
		}
	})

	t.Run("CliSpecs", func(t *testing.T) {
		spec := NewNetworkSpec()
		actual := spec.CliSpecs
		expected := []CliSpec{
			CliSpec{Name: "cli.peer0.org1.fabric.network", ChannelNames: []string{"ChannelOne", "ChannelTwo"}, OrdererDomainName: "orderer.fabric.network"},
			CliSpec{Name: "cli.peer0.org2.fabric.network", ChannelNames: []string{"ChannelOne", "ChannelTwo"}, OrdererDomainName: "orderer.fabric.network"},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected: %v Got: %v", expected, actual)
		}
	})

}

// Generate crypto spec test
func tfixtureCryptoConfigYAMLExists(t *testing.T) func() {

	t.Helper()

	if _, err := os.Stat("crypto-config.yaml"); os.IsNotExist(err) {
		t.Fatalf("crypto-config.yaml does not exists: %v", err)
	}

	return func() { os.Remove("crypto-config.yaml") }

}

func tfixtureVerifyCryptoYAMLFormatting(t *testing.T) {

	t.Helper()

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("crypto-config")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("crypto-config.yaml error %v", err)
	}

}

func TestGenerateCryptoSpec(t *testing.T) {

	data := NetworkSpec{}
	data.NetworkPath = "."
	data.Domain = "fabric.network"
	data.ConsortiumSpecs = []ConsortiumSpec{
		{
			Name: "SampleConsortium",
			ChannelSpecs: []ChannelSpec{
				ChannelSpec{Name: "ChannelOne", Organizations: []string{"Org1", "Org2"}},
				ChannelSpec{Name: "ChannelTwo", Organizations: []string{"Org2"}},
			},
		},
	}
	data.OrganizationSpecs = []OrgSpec{
		{
			Name:   "Org1",
			ID:     "Org1MSP",
			Anchor: "peer0",
		},
		{
			Name:   "Org2",
			ID:     "Org2MSP",
			Anchor: "peer0",
		},
		{
			Name:   "Org3",
			ID:     "Org3MSP",
			Anchor: "peer0",
		},
	}

	generateCryptoSpec(data)
	cleanup := tfixtureCryptoConfigYAMLExists(t)
	defer cleanup()
	tfixtureVerifyCryptoYAMLFormatting(t)
}

// Generate configtx spec test
func tfixtureConfigtxYAMLExists(t *testing.T) func() {

	t.Helper()

	if _, err := os.Stat("configtx.yaml"); os.IsNotExist(err) {
		t.Fatalf("configtx.yaml does not exists: %v", err)
	}

	return func() { os.Remove("configtx.yaml") }

}

func tfixtureVerifyConfigtxYAMLFormatting(t *testing.T) {

	t.Helper()

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("configtx")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("Configtx file error %v", err)
	}

}

func TestGenerateConfigtxSpec(t *testing.T) {

	data := NetworkSpec{}
	data.NetworkPath = "."
	data.Domain = "fabric.network"
	data.ConsortiumSpecs = []ConsortiumSpec{
		{
			Name: "SampleConsortium",
			ChannelSpecs: []ChannelSpec{
				ChannelSpec{Name: "ChannelOne", Organizations: []string{"Org1", "Org2"}},
				ChannelSpec{Name: "ChannelTwo", Organizations: []string{"Org2"}},
			},
		},
	}
	data.OrganizationSpecs = []OrgSpec{
		{
			Name:   "Org1",
			ID:     "Org1MSP",
			Anchor: "peer0",
		},
		{
			Name:   "Org2",
			ID:     "Org2MSP",
			Anchor: "peer0",
		},
		{
			Name:   "Org3",
			ID:     "Org3MSP",
			Anchor: "peer0",
		},
	}

	generateConfigTxSpec(data)
	cleanup := tfixtureConfigtxYAMLExists(t)
	defer cleanup()
	tfixtureVerifyConfigtxYAMLFormatting(t)
}
