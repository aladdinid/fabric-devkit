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
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const cryptoSpecTemplateText = `
# ---------------------------------------------------------------------------
# "OrdererOrgs" - Definition of organizations managing orderer nodes
# ---------------------------------------------------------------------------
OrdererOrgs:
  # ---------------------------------------------------------------------------
  # Orderer
  # ---------------------------------------------------------------------------
  - Name: Orderer
    Domain: {{.Domain}}
    # ---------------------------------------------------------------------------
    # "Specs" - See PeerOrgs below for complete description
    # ---------------------------------------------------------------------------
    Specs:
      - Hostname: orderer

# ---------------------------------------------------------------------------
# "PeerOrgs" - Definition of organizations managing peer nodes
# ---------------------------------------------------------------------------
PeerOrgs:
{{- $domain := .Domain}}
{{- range $index, $org := .OrganizationSpecs }}
  # ---------------------------------------------------------------------------
  # {{$org.Name}}
  # ---------------------------------------------------------------------------
  - Name: {{$org.Name}}
    Domain: {{$org.Name | ToLower}}.{{$domain}}
    # ---------------------------------------------------------------------------
    # "Specs"
    # ---------------------------------------------------------------------------
    # Uncomment this section to enable the explicit definition of hosts in your
    # configuration.  Most users will want to use Template, below
    #
    # Specs is an array of Spec entries.  Each Spec entry consists of two fields:
    #   - Hostname:   (Required) The desired hostname, sans the domain.
    #   - CommonName: (Optional) Specifies the template or explicit override for
    #                 the CN.  By default, this is the template:
    #
    #                              "{Hostname}.{Domain}"
    #
    #                 which obtains its values from the Spec.Hostname and
    #                 Org.Domain, respectively.
    # ---------------------------------------------------------------------------
    # Specs:
    #   - Hostname: foo # implicitly "foo.org1.example.com"
    #     CommonName: foo27.org5.example.com # overrides Hostname-based FQDN set above
    #   - Hostname: bar
    #   - Hostname: baz
    # ---------------------------------------------------------------------------
    # "Template"
    # ---------------------------------------------------------------------------
    # Allows for the definition of 1 or more hosts that are created sequentially
    # from a template. By default, this looks like "peer%d" from 0 to Count-1.
    # You may override the number of nodes (Count), the starting index (Start)
    # or the template used to construct the name (Hostname).
    #
    # Note: Template and Specs are not mutually exclusive.  You may define both
    # sections and the aggregate nodes will be created for you.  Take care with
    # name collisions
    # ---------------------------------------------------------------------------
    Template:
      Count: 1
      # Start: 5
      # Hostname: {.Prefix}{.Index} # default
    # ---------------------------------------------------------------------------
    # "Users"
    # ---------------------------------------------------------------------------
    # Count: The number of user accounts _in addition_ to Admin
    # ---------------------------------------------------------------------------
    Users:
       Count: 1
{{- end }}
`

func generateCryptoSpec(spec NetworkSpec) error {

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tpl := template.Must(template.New("Main").Funcs(funcMap).Parse(cryptoSpecTemplateText))
	configtxYml := filepath.Join(spec.NetworkPath, "crypto-config.yaml")
	f, err := os.Create(configtxYml)
	if err != nil {
		return err
	}
	err = tpl.Execute(f, spec)
	if err != nil {
		return err
	}

	return nil
}

const cryptoConfigExecScriptText = `#!/bin/bash
cryptogen generate --config=./crypto-config.yaml --output="./crypto-config"
`

func generateCryptoExecScript(spec NetworkSpec) error {

	scriptBody := []byte(cryptoConfigExecScriptText)
	generateCryptoExecSh := filepath.Join(spec.NetworkPath, "generateCryptoAsset.sh")
	err := ioutil.WriteFile(generateCryptoExecSh, scriptBody, 0777)
	if err != nil {
		return err
	}

	return nil
}

// GenerateCryptoAssests produces crypo assets and place it in location defined by network spec
func execCryptoScript(spec NetworkSpec) error {

	cmd := []string{"./generateCryptoAsset.sh"}
	if err := RunCryptoConfigContainer(spec.NetworkPath, "cryptogen", "hyperledger/fabric-tools", cmd); err != nil {
		return err
	}
	return nil
}

// CreateCryptoArtifacts produce crypto artefacts
func CreateCryptoArtifacts(spec NetworkSpec) error {
	if err := generateCryptoSpec(spec); err != nil {
		return err
	}

	if err := generateCryptoExecScript(spec); err != nil {
		return err
	}

	if err := execCryptoScript(spec); err != nil {
		return err
	}
	return nil
}
