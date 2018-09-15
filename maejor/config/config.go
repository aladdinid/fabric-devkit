package config

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// ConfigFilename is the default name for the configuration file
const ConfigFilename = ".maejor.yaml"

var configTemplate = template.Must(template.New(".maejor.yaml").Parse(`ProjectPath: {{.ProjectPath}}
ConsortiumPath: {{.ProjectPath}}/network
CryptoConfigPath: {{.ProjectPath}}/crypto-config
ChannelConfigPath: {{.ProjectPath}}/channel-artefacts

network:
   domain: "fabric.network"
   organizations:
     - org1
     - org2
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
func Initialize() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.AddConfigPath(pwd)
	viper.SetConfigName(".maejor")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
