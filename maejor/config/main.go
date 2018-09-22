// +build ignore

package main

import (
	"log"
	"os"

	"github.com/aladdinid/fabric-devkit/maejor/config"
)

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configFiles := config.Search(pwd)
	if len(configFiles) != 0 {
		os.Remove(configFiles[0])
	}

	if err := config.Create(pwd, pwd); err != nil {
		log.Fatal(err)
	}
}
