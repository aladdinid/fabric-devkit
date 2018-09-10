package config

import (
	"os"
	"path/filepath"
	"text/template"
)

var configTemplate = template.Must(template.New(".maejor.yaml").Parse(`
ProjectPath: {{.ProjectPath}}
ConsortiumPath: {{.ProjectPath}}/network
CryptoConfigPath: {{.ProjectPath}}/crypto-config
ChannelConfigPath: {{.ProjectPath}}/channel-artefacts
`))

// Create create configuration file ".maejor.yaml" in
// specified location is one does not exists
func Create(configPath string, projectPath string) error {
	configFile := filepath.Join(configPath, ".maejor.yaml")
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
