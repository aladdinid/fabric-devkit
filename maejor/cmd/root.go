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

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aladdinid/fabric-devkit/internal/config"
	"github.com/spf13/cobra"
)

const description = `maejor is the command-line interface (cli) for 
Aladdin's Hyperledger Fabric (Fabric) development kit.

Copyright 2018 Aladdin Blockchain Technologies Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

var networkSpec *config.NetworkSpec

var rootCmd = &cobra.Command{
	Use:   "maejor",
	Short: "maejor is a cli for Aladdin's Fabric Developer Kit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(description)
		createNetworkPath()
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	cobra.OnInitialize(initConfig)
}

func createNetworkPath() {

	networkPath := config.NetworkPath()
	if err := os.MkdirAll(networkPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := config.Create(pwd, pwd); err != nil {
		log.Fatal(err)
	}

	if err := config.Initialize(pwd, config.ConfigName); err != nil {
		log.Fatal(err)
	}

	networkSpec = config.NewNetworkSpec()
}

// Execute cobra chain of commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
