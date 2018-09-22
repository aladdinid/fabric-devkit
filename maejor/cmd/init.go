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

	"github.com/aladdinid/fabric-devkit/maejor/config"
	"github.com/aladdinid/fabric-devkit/maejor/docker"
	"github.com/spf13/cobra"
)

var imagesPulled bool
var imagesRemoved bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init initialize the project",
	Run: func(cmd *cobra.Command, args []string) {
		hyperledger := config.HyperledgerImages()
		if imagesPulled {
			pullAndRetagImages(hyperledger)
		}
		if imagesRemoved {
			err := deleteImages(hyperledger)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	initCmd.Flags().BoolVarP(&imagesPulled, "pull", "p", false, "pull project images from docker hub")
	initCmd.Flags().BoolVarP(&imagesRemoved, "remove", "r", false, "remove project images")
}

func pullAndRetagImages(images []string) {
	docker.PullImages(images)
	docker.TagImagesAsLatest(images)
}

func deleteImages(images []string) error {
	err := docker.DeleteImages(images)
	if err != nil {
		return err
	}
	return nil
}
