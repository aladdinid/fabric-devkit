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
	"os"
	"testing"

	"github.com/spf13/viper"
)

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
	data.OrganizationSpecs = []OrgSpec{
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
		OrgSpec{
			Name:   "Org3",
			ID:     "Org3MSP",
			Anchor: "peer0",
		},
	}
	data.ConsortiumSpecs = []ConsortiumSpec{
		{
			Name:          "SampleConsortium",
			ChannelName:   "TwoOrg",
			Organizations: []string{"Org1", "Org2"},
		},
	}

	generateCryptoSpec(data)
	cleanup := tfixtureCryptoConfigYAMLExists(t)
	defer cleanup()
	tfixtureVerifyCryptoYAMLFormatting(t)
}
