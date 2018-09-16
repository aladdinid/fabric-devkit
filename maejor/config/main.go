// +build ignore

package main

import (
	"os"

	"github.com/aladdinid/fabric-devkit/maejor/config"
)

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configFiles := config.Search(pwd)
	os.Remove(configFiles[0])

	if err := config.Create(pwd, pwd); err != nil {
		panic(err)
	}
}
