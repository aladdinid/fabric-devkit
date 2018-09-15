package config

import (
	"github.com/spf13/viper"
)

// Hyperledger return a list of docker images
func Hyperledger() []string {
	return viper.GetStringSlice("containers.hyperledger")
}
