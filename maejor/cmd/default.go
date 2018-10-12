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
	"log"

	"github.com/aladdinid/fabric-devkit/internal/config"
	"github.com/aladdinid/fabric-devkit/internal/configtx"
	"github.com/aladdinid/fabric-devkit/internal/crypto"
	"github.com/aladdinid/fabric-devkit/internal/network"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Create a default network specification",
	Run: func(cmd *cobra.Command, args []string) {
		hyperledger := config.HyperledgerImages()
		pullAndRetagImages(hyperledger)
		createNetworkPath()
		createChannelArtefactPath()
		if err := configtx.GenerateConfigtxSpec(*networkSpec); err != nil {
			log.Fatal(err)
		}
		if err := crypto.GenerateCryptoSpec(*networkSpec); err != nil {
			log.Fatal(err)
		}
		if err := network.GenerateNetworkSpec(*networkSpec); err != nil {
			log.Fatal(err)
		}
	},
}
