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

var imagesRemoved bool

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Manage relevant docker images for the project",
	Run: func(cmd *cobra.Command, args []string) {
		hyperledger := svc.HyperledgerImages()
		pullAndRetagImages(hyperledger)
		if imagesRemoved {
			err := deleteImages(hyperledger)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	imageCmd.Flags().BoolVarP(&imagesRemoved, "reset", "r", false, "delete and pull images")
}

func pullAndRetagImages(images []string) {
	if err := svc.PullImages(images); err != nil {
		log.Fatal(err)
	}
	if err := svc.TagImages(images, svc.TargetTagAsLatest); err != nil {
		log.Fatal(err)
	}
}

func deleteImages(images []string) error {
	err := svc.DeleteImages(images)
	if err != nil {
		return err
	}
	return nil
}
