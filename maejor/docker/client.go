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
	"os"
	"strings"

	"github.com/docker/docker/api/types"
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

// PullImages pull multiple images
func PullImages(names []string) {

	for _, name := range names {
		reader, err := pullImage(name)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, reader)
		reader.Close()
	}
}

// TagImages multiple sources and targets
func TagImages(sources []string, targets []string) {

	if len(sources) != len(targets) {
		log.Fatal(fmt.Errorf("sources and targets did not match"))
	}

	for index, source := range sources {
		err := tagImage(source, targets[index])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// TagImagesAsLatest all source images to contain "latest" suffix
func TagImagesAsLatest(sources []string) {

	targets := tagImagesAsLatest(sources)
	TagImages(sources, targets)

}
