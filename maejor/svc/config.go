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
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// ConfigName name of the config file name without extension
const ConfigName = ".maejor"

const configTemplateText = `
ProjectPath: {{.ProjectPath}}
NetworkPath: {{.ProjectPath}}/network
CryptoPath: {{.ProjectPath}}/network/crypto-config
ChannelArtefactPath: {{.ProjectPath}}/network/channel-artefacts
ScriptPath: {{.ProjectPath}}/network/scripts
ChaincodePath: $GOPATH/src/github.com/aladdinid/chaincodes

containers:
  hyperledger:
     - hyperledger/fabric-ca:x86_64-1.1.0
     - hyperledger/fabric-peer:x86_64-1.1.0
     - hyperledger/fabric-orderer:x86_64-1.1.0
     - hyperledger/fabric-ccenv:x86_64-1.1.0
     - hyperledger/fabric-tools:x86_64-1.1.0
     - hyperledger/fabric-couchdb:x86_64-1.0.6

network:
   domain: "fabric.network"
   consortiums: 
     - SampleConsortium
   organizations:
     - Org1
     - Org2

SampleConsortium:
   name: "SampleConsortium"
   channelname: "TwoOrg"
   organizations:
     - Org1
     - Org2

Org1:
  name: Org1
  id: Org1MSP
  anchor: peer0

Org2:
  name: Org2
  id: Org2MSP
  anchor: peer0
`

var configFilename = strings.Join([]string{ConfigName, ".yaml"}, "")

// ConfigTemplate contain template engines
var ConfigTemplate = template.Must(template.New(".maejor.yaml").Parse(configTemplateText))

// Create create configuration file ".maejor.yaml" in
// specified location is one does not exists
func Create(configPath string, projectPath string) error {
	configFile := filepath.Join(configPath, configFilename)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err := os.Create(configFile)
		if err != nil {
			return err
		}

		err = ConfigTemplate.Execute(f, struct {
			ProjectPath string
		}{
			projectPath,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Search for configuration file
func Search(rootPath string) []string {

	result := []string{}
	err := filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, configFilename) {
			result = append(result, path)
		}
		return nil
	})

	if err != nil {
		return []string{}
	}

	return result
}

// InitializeByReader sets viper engine by byte slice
// This is intended primarily to support testing
func InitializeByReader(config []byte) error {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(config)); err != nil {
		return err
	}
	return nil
}

// Initialize sets viper via configuration file.
// This is the main entry point.
func Initialize(configPath string, configName string) error {

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// ProjectPath returns the project path
func ProjectPath() string {
	return viper.GetString("ProjectPath")
}

// NetworkPath returns path to the network spec
func NetworkPath() string {
	return viper.GetString("NetworkPath")
}

// CryptoPath returns path to crypto assets
func CryptoPath() string {
	return viper.GetString("CryptoPath")
}

// ChannelArtefactPath returns path to channel artefacts
func ChannelArtefactPath() string {
	return viper.GetString("ChannelArtefactPath")
}

// ScriptPath return path to scripts
func ScriptPath() string {
	return viper.GetString("ScriptPath")
}

// ChaincodePath returns path to chaincode
func ChaincodePath() string {
	return viper.GetString("ChaincodePath")
}

// HyperledgerImages return a list of fabric images
func HyperledgerImages() []string {
	return viper.GetStringSlice("containers.hyperledger")
}

// Domain returns configuration
func Domain() string {
	return viper.GetString("network.domain")
}

func consortiumByName(name string) ConsortiumSpec {
	value := viper.GetStringMap(name)

	var spec = ConsortiumSpec{}

	name, ok := value["name"].(string)
	if ok {
		spec.Name = name
	}

	channelName, ok := value["channelname"].(string)
	if ok {
		spec.ChannelName = channelName
	}

	iOrgs, ok := value["organizations"].([]interface{})
	if ok {
		var orgs []string
		for _, iOrg := range iOrgs {
			switch s := iOrg.(type) {
			case string:
				orgs = append(orgs, s)
			}
		}
		spec.Organizations = orgs

	}

	return spec
}

// ConsortiumSpecs returns am array of consortium specs
func ConsortiumSpecs() []ConsortiumSpec {

	var consortiumSpecs = []ConsortiumSpec{}

	consortiumNames := viper.GetStringSlice("network.consortiums")

	for _, consortiumName := range consortiumNames {
		consortiumSpec := consortiumByName(consortiumName)
		consortiumSpecs = append(consortiumSpecs, consortiumSpec)
	}

	return consortiumSpecs
}

func orgByName(name string) OrgSpec {
	value := viper.GetStringMap(name)

	var spec = OrgSpec{}

	name, ok := value["name"].(string)
	if ok {
		spec.Name = name
	}

	id, ok := value["id"].(string)
	if ok {
		spec.ID = id
	}

	anchor, ok := value["anchor"].(string)
	if ok {
		spec.Anchor = anchor
	}

	return spec
}

// OrganizationSpecs returns an array of organizations specs
func OrganizationSpecs() []OrgSpec {

	var orgSpecs = []OrgSpec{}

	orgNames := viper.GetStringSlice("network.organizations")
	for _, orgName := range orgNames {
		orgSpec := orgByName(orgName)
		orgSpecs = append(orgSpecs, orgSpec)
	}

	return orgSpecs
}
