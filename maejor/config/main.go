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

	if err := config.Create(pwd, pwd); err != nil {
		panic(err)
	}
}
