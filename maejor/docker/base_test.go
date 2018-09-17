package docker

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

func pullTestImages(t *testing.T) {

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

	pullTestImages(t)

	result, err := searchImages("alpine:*")
	if err != nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected: 2 Got %d", len(result))
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
