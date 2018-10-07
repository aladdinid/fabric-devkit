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

package configtx_test

import (
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/aladdinid/fabric-devkit/internal/config"
	"github.com/aladdinid/fabric-devkit/internal/configtx"
)

func fixtureConfigtxYAMLExists(t *testing.T) func() {

	t.Helper()

	if _, err := os.Stat("configtx.yaml"); os.IsNotExist(err) {
		t.Fatalf("configtx.yaml does not exists: %v", err)
	}

	return func() { os.Remove("configtx.yaml") }

}

func fixtureVerifyConfigtxYAMLFormatting(t *testing.T) {

	t.Helper()

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("configtx")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("Configtx file error %v", err)
	}

}

func TestGenerateConfigtxSpec(t *testing.T) {

	data := config.NetworkSpec{}
	data.NetworkPath = "."
	data.Domain = "fabric.network"
	data.OrganisationSpecs = []config.OrgSpec{
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
		config.OrgSpec{
			Name:   "Org3",
			ID:     "Org3MSP",
			Anchor: "peer0",
		},
	}
	data.ChannelName = "Test"
	data.Consortium = "SampleConsortium"

	configtx.GenerateConfigtxSpec(data)
	cleanup := fixtureConfigtxYAMLExists(t)
	defer cleanup()
	fixtureVerifyConfigtxYAMLFormatting(t)
}
