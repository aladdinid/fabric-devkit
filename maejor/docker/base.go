/*
Copyright 2018 Paul Sitoh
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

func pullImage(name string) (io.ReadCloser, error) {

	if cli == nil {
		return nil, fmt.Errorf("session not started")
	}

	reader, err := cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func tagImage(source string, target string) error {

	if cli == nil {
		return fmt.Errorf("session not started")
	}

	err := cli.ImageTag(ctx, source, target)
	if err != nil {
		return err
	}

	return nil
}

func tagImageAsLatest(source string) string {

	result := strings.Split(source, ":")
	result[1] = "latest"
	return strings.Join(result, ":")

}

func tagImagesAsLatest(sources []string) []string {

	var result []string
	for _, source := range sources {
		replacement := tagImageAsLatest(source)
		result = append(result, replacement)
	}

	return result
}

func searchImages(source string) ([]string, error) {

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

// DeletedImage is a structure that represents
type DeletedImage struct {
	Deleted  string
	Untagged string
}

func removeImage(imageID string) ([]types.ImageDeleteResponseItem, error) {

	if cli == nil {
		return []types.ImageDeleteResponseItem{}, fmt.Errorf("session not started")
	}

	deletes, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{})
	if err != nil {
		return []types.ImageDeleteResponseItem{}, err
	}

	return deletes, nil
}
