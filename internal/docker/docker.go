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

package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var cli *client.Client
var ctx = context.Background()

func init() {
	c, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c
}

// PullImage initate a pull from Dockerhub
func PullImage(name string) (io.ReadCloser, error) {

	if cli == nil {
		return nil, fmt.Errorf("session not started")
	}

	reader, err := cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// PullImages pull multiple images
func PullImages(names []string) {

	for _, name := range names {
		reader, err := PullImage(name)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, reader)
		reader.Close()
	}
}

// TagImage downloaded images
func TagImage(source string, target string) error {

	if cli == nil {
		return fmt.Errorf("session not started")
	}

	err := cli.ImageTag(ctx, source, target)
	if err != nil {
		return err
	}

	return nil
}

// TagImages multiple sources and targets
func TagImages(sources []string, targets []string) {

	if len(sources) != len(targets) {
		log.Fatal(fmt.Errorf("sources and targets did not match"))
	}

	for index, source := range sources {
		err := TagImage(source, targets[index])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// TagImageAsLatest tag a specific image as latest
func TagImageAsLatest(source string) string {

	result := strings.Split(source, ":")
	result[1] = "latest"
	return strings.Join(result, ":")

}

// TagImagesAsLatest tag a list of images as latest
func TagImagesAsLatest(sources []string) []string {

	var result []string
	for _, source := range sources {
		replacement := TagImageAsLatest(source)
		result = append(result, replacement)
	}

	return result
}

// SearchImages search for images by name
func SearchImages(source string) ([]string, error) {

	if cli == nil {
		return nil, fmt.Errorf("session not started")
	}

	filters := filters.NewArgs()
	filters.Add("reference", source)

	images, _ := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters,
	})

	ids := []string{}
	for _, image := range images {
		ids = append(ids, image.ID)
	}

	return ids, nil
}

// RemoveImage by image id
func RemoveImage(imageID string) ([]types.ImageDeleteResponseItem, error) {

	if cli == nil {
		return []types.ImageDeleteResponseItem{}, fmt.Errorf("session not started")
	}

	deletes, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force: true,
	})
	if err != nil {
		return []types.ImageDeleteResponseItem{}, err
	}

	return deletes, nil
}

// DeleteImages remove downloaded images
func DeleteImages(images []string) error {
	for _, image := range images {
		ids, err := SearchImages(image)
		if err != nil {
			return err
		}

		for _, id := range ids {
			result, err := RemoveImage(id)
			log.Printf("Images removed %v", result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
