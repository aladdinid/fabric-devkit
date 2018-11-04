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
	"strings"

	"github.com/aladdinid/fabric-devkit/maejor/svc"
	"github.com/spf13/cobra"
)

var generateNetworkSpecCmd = &cobra.Command{
	Use:   "generate",
	Short: "Skeleton network specification",
	Run: func(cmd *cobra.Command, args []string) {
		if err := createNetworkSpec(); err != nil {
			if strings.Contains(err.Error(), "Error: No such image: hyperledger/fabric-tools") {
				fmt.Println(err)
				fmt.Println(`Please run the command: "maejor container image pull"`)
			} else {
				log.Fatal(err)
			}
		}
	},
}

var deleteNetworkSpecCmd = &cobra.Command{
	Use:   "delete",
	Short: "Network specification",
	Run: func(cmd *cobra.Command, args []string) {
		var yesResponse string

		fmt.Print("Delete you project specs [N/y]? ")
		fmt.Scanln(&yesResponse)
		if strings.Compare(strings.ToLower(yesResponse), "y") == 0 || strings.Compare(strings.ToLower(yesResponse), "yes") == 0 {
			if err := deleteNetworkSpec(); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Operations to generate network specification",
}

func init() {
	networkCmd.AddCommand(generateNetworkSpecCmd)
	networkCmd.AddCommand(deleteNetworkSpecCmd)
}

func createNetworkSpec() error {
	networkPath := svc.NetworkPath()
	if err := os.MkdirAll(networkPath, os.ModePerm); err != nil {
		return err
	}

	channelArtefactPath := svc.ChannelArtefactPath()
	if err := os.MkdirAll(channelArtefactPath, os.ModePerm); err != nil {
		return err
	}

	scriptPath := svc.ScriptPath()
	if err := os.MkdirAll(scriptPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// You must create crypto first then channel artefacts
	if err := svc.CreateCryptoArtifacts(*networkSpec); err != nil {
		return err
	}
	if err := svc.CreateChannelArtefacts(*networkSpec); err != nil {
		return err
	}
	if err := svc.CreateNetworkSpec(*networkSpec); err != nil {
		return err
	}

	if err := svc.GenerateCliScripts(*networkSpec); err != nil {
		return err
	}

	return nil
}

func deleteNetworkSpec() error {

	networkPath := svc.NetworkPath()

	if err := os.RemoveAll(networkPath); err != nil {
		return err
	}

	return nil
}
