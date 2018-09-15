package config

import (
	"github.com/spf13/viper"
)

// ProjectPath returns the project path
func ProjectPath() string {
	return viper.GetString("ProjectPath")
}
