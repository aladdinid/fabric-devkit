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

package svc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
func PullImages(names []string) error {

	for _, name := range names {
		reader, err := PullImage(name)
		if err != nil {
			reader.Close()
			return err
		}
		io.Copy(os.Stdout, reader)
		reader.Close()
	}

	return nil
}

// TagImage downloaded images
func TagImage(source string, target func(string) string) error {

	if cli == nil {
		return fmt.Errorf("session not started")
	}

	err := cli.ImageTag(ctx, source, target(source))
	if err != nil {
		return err
	}

	return nil
}

// TagImages multiple sources and targets
func TagImages(sources []string, target func(string) string) error {

	for _, source := range sources {
		err := TagImage(source, target)
		if err != nil {
			return err
		}
	}

	return nil
}

// TargetTagAsLatest ensure that source are tagged as latest
func TargetTagAsLatest(source string) string {

	result := strings.Split(source, ":")
	if len(result) == 2 {
		result[1] = "latest"
	}
	return strings.Join(result, ":")
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

// RunCryptoConfigContainer run a container
func RunCryptoConfigContainer(hostVol, name, image string, cmd []string) error {

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: image,
			Env: []string{
				"GOPATH=/opt/gopath",
				"FABRIC_CFG_PATH=/opt/gopath/src/github.com/hyperledger/fabric",
			},
			WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric",
			Cmd:        cmd,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: hostVol,
					Target: "/opt/gopath/src/github.com/hyperledger/fabric",
				},
			},
		},
		&network.NetworkingConfig{},
		name,
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}

	return nil
}
