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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const description = `maejor is the command-line interface (cli) for 
Aladdin's Hyperledger Fabric (Fabric) development kit.

Copyright 2018 Aladdin Blockchain Technologies Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

var rootCmd = &cobra.Command{
	Use:   "maejor",
	Short: "maejor is a cli for Aladdin's Fabric Developer Kit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(description)
	},
}

// Execute cobra chain of commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
