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
	"fmt"
	"io"
	"log"
	"os"
)

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
