// +build smoke

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

package docker_test

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aladdinid/fabric-devkit/internal/docker"
)

func TestPullImage(t *testing.T) {

	_, err := docker.PullImage("unbuntu")
	if err == nil {
		t.Fatalf("Expect: error Got: no error")
	}

	_, err = docker.PullImage("ubuntu")
	if err != nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}

}

func fixturesForSearchImageTest(t *testing.T) {

	reader, err := docker.PullImage("alpine:latest")
	if err != nil {
		t.Fatal("Unable to pull alpine:latest")
	}

	io.Copy(os.Stdout, reader)

	reader, err = docker.PullImage("alpine:3.7")
	if err != nil {
		t.Fatal("Unable to pull alpine:3.7")
	}

	io.Copy(os.Stdout, reader)
}

func TestSearchImages(t *testing.T) {

	fixturesForSearchImageTest(t)

	result, err := docker.SearchImages("alpine:*")
	if err != nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected: 2 Got %d", len(result))
	}

}

func fixturesForRemoveImageTest(t *testing.T) []string {

	reader, err := docker.PullImage("alpine:3.5")
	if err != nil {
		t.Fatal("Unable to pull alpine:3.5")
	}

	io.Copy(os.Stdout, reader)

	ids, err := docker.SearchImages("alpine:3.5")
	if err != nil {
		t.Fatal("image is not found")
	}

	return ids

}

func TestRemoveImage(t *testing.T) {

	ids := fixturesForRemoveImageTest(t)

	deleted, err := docker.RemoveImage(ids[0])
	if err != nil {
		t.Fatalf("Expected: no err Got: %v", err)
	}

	if len(deleted) != 4 {
		t.Fatalf("Expected: 1 Got: %d", len(deleted))
	}

}

func TestTagImage(t *testing.T) {

	err := docker.TagImage("something", "something else")
	if err == nil {
		t.Fatal("Expected: error Got: no error")
	}

}

func TestTagImageAsLatest(t *testing.T) {

	source := "something:1234"
	expected := "something:latest"
	result := docker.TagImageAsLatest(source)

	if strings.Compare(expected, result) != 0 {
		t.Fatalf("Source: %s Expected: %s Got: %s", source, expected, result)
	}

}

func TestTagImagesAsLatest(t *testing.T) {

	source := []string{"something:1234", "else:1234"}
	expected := []string{"something:latest", "else:latest"}

	result := docker.TagImagesAsLatest(source)
	if reflect.DeepEqual(expected, result) != true {
		t.Fatalf("Source: %v Expected: %v Got: %v", source, expected, result)
	}

}

func fixtureLocation(t *testing.T) string {

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return pwd

}

func fixturePullFabricTools(t *testing.T) {
	reader, err := docker.PullImage("hyperledger/fabric-tools:x86_64-1.1.0")
	if err != nil {
		t.Fatal("Unable to pull fabric-tools:x86_64-1.1.0")
	}

	io.Copy(os.Stdout, reader)

}

func fixtureTagFabricTools(t *testing.T) {
	latest := docker.TagImageAsLatest("hyperledger/fabric-tools:x86_64-1.1.0")
	if strings.Compare(latest, "hyperledger/fabric-tools:latest") != 0 {
		t.Fatalf("Expected: hyperledger/fabric-tools:latest Got: %s", latest)
	}
}

func TestRunCryptoConfigContainer(t *testing.T) {
	fixturePullFabricTools(t)
	fixtureTagFabricTools(t)
	err := docker.RunCryptoConfigContainer(fixtureLocation(t), "/opt/gopath/src/github.com/hyperledger/fabric", "hyperledger/fabric-tools", []string{"which", "cryptogen"})
	if err != nil {
		t.Fatal(err)
	}

	err = docker.RunCryptoConfigContainer(fixtureLocation(t), "/opt/gopath/src/github.com/hyperledger/fabric", "hyperledger/fabric-tools", []string{"which", "configtxgen"})
	if err != nil {
		t.Fatal(err)
	}
}
