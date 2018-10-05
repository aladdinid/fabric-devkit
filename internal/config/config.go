package config

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// ConfigName name of the config file name without extension
const ConfigName = ".maejor"

// ConfigFilename is the default name for the configuration file
var ConfigFilename = strings.Join([]string{ConfigName, ".yaml"}, "")

var configTemplate = template.Must(template.New(".maejor.yaml").Parse(`ProjectPath: {{.ProjectPath}}
ConsortiumPath: {{.ProjectPath}}/network
CryptoConfigPath: {{.ProjectPath}}/network/crypto-config
ChannelConfigPath: {{.ProjectPath}}/network/channel-artefacts

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
   orderer:
       - OrdererOrg
   organizations:
       - Org1
       - Org2
`))

// Create create configuration file ".maejor.yaml" in
// specified location is one does not exists
func Create(configPath string, projectPath string) error {
	configFile := filepath.Join(configPath, ConfigFilename)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err := os.Create(configFile)
		if err != nil {
			return err
		}

		err = configTemplate.Execute(f, struct {
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
		if strings.Contains(path, ConfigFilename) {
			result = append(result, path)
		}
		return nil
	})

	if err != nil {
		return []string{}
	}

	return result
}

// Initialize viper framework
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

// HyperledgerImages return a list of fabric images
func HyperledgerImages() []string {
	return viper.GetStringSlice("containers.hyperledger")
}

// Domain returns configuration
func Domain() string {
	return viper.GetString("network.domain")
}

// Organizations returns
func Organizations() []string {
	return viper.GetStringSlice("network.organizations")
}
