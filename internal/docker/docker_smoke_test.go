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

// +build smoke

package docker_test

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestPullImage(t *testing.T) {

	_, err := pullImage("unbuntu")
	if err == nil {
		t.Fatalf("Expect: error Got: no error")
	}

	_, err = pullImage("ubuntu")
	if err != nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}

}

func fixturesForSearchImageTest(t *testing.T) {

	reader, err := pullImage("alpine:latest")
	if err != nil {
		t.Fatal("Unable to pull alpine:latest")
	}

	io.Copy(os.Stdout, reader)

	reader, err = pullImage("alpine:3.7")
	if err != nil {
		t.Fatal("Unable to pull alpine:3.7")
	}

	io.Copy(os.Stdout, reader)
}

func TestSearchImages(t *testing.T) {

	fixturesForSearchImageTest(t)

	result, err := searchImages("alpine:*")
	if err != nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected: 2 Got %d", len(result))
	}

}

func fixturesForRemoveImageTest(t *testing.T) []string {

	reader, err := pullImage("alpine:3.5")
	if err != nil {
		t.Fatal("Unable to pull alpine:3.5")
	}

	io.Copy(os.Stdout, reader)

	ids, err := searchImages("alpine:3.5")
	if err != nil {
		t.Fatal("image is not found")
	}

	return ids

}

func TestRemoveImage(t *testing.T) {

	ids := fixturesForRemoveImageTest(t)

	deleted, err := removeImage(ids[0])
	if err != nil {
		t.Fatalf("Expected: no err Got: %v", err)
	}

	if len(deleted) != 4 {
		t.Fatalf("Expected: 1 Got: %d", len(deleted))
	}

}

func TestTagImage(t *testing.T) {

	err := tagImage("something", "something else")
	if err == nil {
		t.Fatal("Expected: error Got: no error")
	}

}

func TestTagImageAsLatest(t *testing.T) {

	source := "something:1234"
	expected := "something:latest"
	result := tagImageAsLatest(source)

	if strings.Compare(expected, result) != 0 {
		t.Fatalf("Source: %s Expected: %s Got: %s", source, expected, result)
	}

}

func TestTagImagesAsLatest(t *testing.T) {

	source := []string{"something:1234", "else:1234"}
	expected := []string{"something:latest", "else:latest"}

	result := tagImagesAsLatest(source)
	if reflect.DeepEqual(expected, result) != true {
		t.Fatalf("Source: %v Expected: %v Got: %v", source, expected, result)
	}

}
