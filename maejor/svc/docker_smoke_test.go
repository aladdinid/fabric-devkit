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

package svc

import (
	"io"
	"os"
	"strings"
	"testing"
)

const tfixtureHelloWorldLinux = "hello-world:linux"
const tfixtureHelloWorldLatest = "hello-world"

func TestBasicOperations(t *testing.T) {

	t.Run("PullImage", func(t *testing.T) {
		reader, err := PullImage(tfixtureHelloWorldLinux)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		io.Copy(os.Stdout, reader)
	})

	t.Run("PullImages", func(t *testing.T) {
		if err := PullImages([]string{tfixtureHelloWorldLatest}); err != nil {
			t.Fatalf("Expected no error. Got: %v", err)
		}
	})

	t.Run("SearchImage", func(t *testing.T) {
		if _, err := SearchImages("hello-world:*"); err != nil {
			t.Fatalf("Expected: no error Got: %v", err)
		}

		result, err := SearchImages("hello-world:*")
		if err != nil {
			t.Fatalf("Expected: no error Got: %v", err)
		}

		if len(result) != 1 {
			t.Fatalf("Expected: 1 Got %d", len(result))
		}
	})

	t.Run("RemoveImage", func(t *testing.T) {
		ids, err := SearchImages("hello-world")
		if err != nil {
			t.Fatal("image is not found")
		}

		deleted, err := RemoveImage(ids[0])
		if err != nil {
			t.Fatalf("Expected: no err Got: %v", err)
		}

		if len(deleted) != 6 {
			t.Fatalf("Expected: 6 Got: %d", len(deleted))
		}
	})

}

func tfixtureRemoveImage(t *testing.T) {
	ids, err := SearchImages("hello-world")
	if err != nil {
		t.Fatal("image is not found")
	}

	_, err = RemoveImage(ids[0])
	if err != nil {
		t.Fatalf("Error not expected. Got: %v", err)
	}
}

func tfixturePullDummyImage(t *testing.T) {
	reader, err := PullImage(tfixtureHelloWorldLinux)
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	io.Copy(os.Stdout, reader)
}

func tfixtureTagTargetAsSomething(source string) string {
	result := strings.Split(source, ":")
	if len(result) == 2 {
		result[1] = "something"
	}
	return strings.Join(result, ":")
}

func TestTaggingImages(t *testing.T) {

	t.Run("TagImage", func(t *testing.T) {
		tfixturePullDummyImage(t)

		err := TagImage(tfixtureHelloWorldLinux, TargetTagAsLatest)
		if err != nil {
			t.Fatalf("Error not expected. Got: %v", err)
		}
	})

	t.Run("TagImages", func(t *testing.T) {

		t.Helper()
		sources := []string{tfixtureHelloWorldLinux, tfixtureHelloWorldLatest}
		if err := TagImages(sources, tfixtureTagTargetAsSomething); err != nil {
			t.Fatalf("Expected no err. Got %v", err)
		}

		tfixtureRemoveImage(t)
	})

}

func tfixtureLocation(t *testing.T) string {

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return pwd

}

func tfixturePullFabricTools(t *testing.T) {
	reader, err := PullImage("hyperledger/fabric-tools:x86_64-1.1.0")
	if err != nil {
		t.Fatal("Unable to pull fabric-tools:x86_64-1.1.0")
	}

	io.Copy(os.Stdout, reader)

}

func tfixtureTagFabricTools(t *testing.T) {
	err := TagImage("hyperledger/fabric-tools:x86_64-1.1.0", TargetTagAsLatest)
	if err != nil {
		t.Fatalf("Expected no error. Got %v", err)
	}
}

func TestRunCryptoConfigContainer(t *testing.T) {
	tfixturePullFabricTools(t)
	tfixtureTagFabricTools(t)
	err := RunCryptoConfigContainer(tfixtureLocation(t), "container_1", "hyperledger/fabric-tools", []string{"which", "cryptogen"})
	if err != nil {
		t.Fatal(err)
	}

	err = RunCryptoConfigContainer(tfixtureLocation(t), "container_2", "hyperledger/fabric-tools", []string{"which", "configtxgen"})
	if err != nil {
		t.Fatal(err)
	}
}
