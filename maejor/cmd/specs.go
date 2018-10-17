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

	"github.com/aladdinid/fabric-devkit/maejor/svc"
	"github.com/spf13/cobra"
)

var specCmd = &cobra.Command{
	Use:   "specs",
	Short: "Create artefacts",
	Run: func(cmd *cobra.Command, args []string) {
		createNetworkPath()
		createChannelArtefactPath()
		if err := svc.GenerateConfigtxSpec(*networkSpec); err != nil {
			log.Fatal(err)
		}
		if err := svc.GenerateConfigTxExecScript(*networkSpec); err != nil {
			log.Fatal(err)
		}
		if err := svc.GenerateCryptoSpec(*networkSpec); err != nil {
			log.Fatal(err)
		}
		if err := svc.GenerateCryptoExecScript(*networkSpec); err != nil {
			log.Fatal(err)
		}

		if err := svc.GenerateCryptoAssests(*networkSpec); err != nil {
			log.Fatal(err)
		}

		if err := svc.GenerateConfigTxAssets(*networkSpec); err != nil {
			log.Fatal(err)
		}

	},
}
