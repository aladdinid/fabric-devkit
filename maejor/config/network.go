package config

import (
	"github.com/spf13/viper"
)

// Domain returns configuration
func Domain() string {
	return viper.GetString("network.domain")
}

// Organizations returns a list of organizations
func Organizations() []string {
	return viper.GetStringSlice("network.organizations")
}
