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
	"io/ioutil"
	"path/filepath"
)

const createChannelScriptTxt = `#/bin/bash

function usage(){
    echo "Usage: $0 <flags>"
    echo "Mandatory:"
    echo "  -c <channelname>  The name of the channel that you wish to create"
    echo "  -o <orderer>      A url to the orderer"
}

while getopts "o:c:l:p:v:" opt; do
  case $opt in
    c)
      CHANNELNAME=$OPTARG
      ;;
    o)
      ORDERER=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

peer channel create -o $ORDERER -c $CHANNELNAME -f ./channel-artefacts/$CHANNELNAME/channel.tx --tls --cafile $ORDERER_CA
`

func generateCreateChannelScript(spec NetworkSpec) error {
	scriptBody := []byte(createChannelScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "create-channel.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}
	return nil

}

const joinChannelScriptTxt = `#/bin/bash

function usage(){
  echo "Usage: $0 <flags>"
  echo "Mandatory:"
  echo "  -c <channelname>  The name of the channel that you wish to create"
  echo "  -o <orderer>      A url to the orderer"
}

while getopts "o:c:l:p:v:" opt; do
  case $opt in
    c)
      CHANNELNAME=$OPTARG
      ;;
    o)
      ORDERER=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

if [ -f "./$CHANNELNAME/$CHANNELNAME.block" ]; then
   peer channel fetch newest -o $ORDERER -c $CHANNELNAME --tls --cafile $ORDERER_CA ./CHANNELNAME/$CHANNELNAME.block
fi

peer channel join -o $ORDERER -b ./CHANNELNAME/$CHANNELNAME.block --tls --cafile $ORDERER_CA
`

func generateJoinChannelScript(spec NetworkSpec) error {
	scriptBody := []byte(joinChannelScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "join-channel.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}
	return nil
}

const installChaincodeScriptTxt = `#!/bin/bash

function usage(){ 
  echo "Usage: $0 <flags>"
  echo "Mandatory:"
  echo "   -c <cc id>         A unique string identifier"
  echo "   -v <cc version>    A numeric number"
  echo "   -p <cc package>    A name of folder containing chaincodes"
  echo "Optional:"
  echo "   -a <cc constructor>  Must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
  echo "   -l <language>        Option to deloy in language other Go, currenly only node"
}

echo "the number of arguments are: $#"
if [ "$#" -eq "0" ]; then  
  usage
  exit
fi

while getopts "a:c:l:p:v:" opt; do
  case $opt in
    a)
      CHAINCODE_CONSTRUCTOR=$OPTARG
      ;;
    c)
      CHAINCODEID=$OPTARG
      ;;
    l)
      CHAINCODE_LANG=$OPTARG
      ;;
    p)
      CHAINCODE_PACKAGE=$OPTARG
      ;;
    v)
      CHAINCODE_VERSION=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

if [ -z $CHAINCODE_LANG ]; then
CHAINCODE_LANG="golang"
fi

if [ -z $CHAINCODE_CONSTRUCTOR ]; then
  CHAINCODE_CONSTRUCTOR="[]"
fi

if [[ ! -z $CHAINCODEID && ! -z $CHAINCODE_VERSION && ! -z $CHAINCODE_PACKAGE ]]; then

  if [ "$CHAINCODE_LANG" == "golang" ]; then 
    path_to_chaincode="github.com/hyperledger/fabric/chaincodes/$CHAINCODE_PACKAGE"
  else
    path_to_chaincode="$GOPATH/src/github.com/hyperledger/fabric/chaincodes/$CHAINCODE_PACKAGE"
  fi
  
  echo "INSTALLING chaincode $CHAINCODEID in $CHAINCODE_LANG version $CHAINCODE_VERSION in $path_to_chaincode"
  echo
  peer chaincode install -n $CHAINCODEID -v $CHAINCODE_VERSION -l $CHAINCODE_LANG -p $path_to_chaincode --tls --cafile $ORDERER_CA

  echo "INSTALLATION chaincode successfully"
else
  usage
fi
`

func generateInstallChaincodeScript(spec NetworkSpec) error {
	scriptBody := []byte(installChaincodeScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "install-chaincode.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}
	return nil
}

const instantiatieChaincodeScriptTxt = `#!/bin/bash

function usage(){ 
  echo "Usage: $0 <flags>"
  echo "Mandatory:"
  echo "   -c <cc id>         A unique string identifier"
  echo "   -v <cc version>    A numeric number"
  echo "   -n <channelname>   A string identifying a channel"
  echo "Optional:"
  echo "   -a <cc constructor>  Must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
  echo "   -l <language>        Option to deloy in language other Go, currenly only node"
}

if [ "$#" -eq "0" ]; then  
  usage
  exit
fi

while getopts "a:c:l:n:v:" opt; do
  case $opt in
    a)
      CHAINCODE_CONSTRUCTOR=$OPTARG
      ;;
    c)
      CHAINCODEID=$OPTARG
      ;;
    l)
      CHAINCODE_LANG=$OPTARG
      ;;
    n)
      CHANNELNAME=$OPTARG
      ;;
    v)
      CHAINCODE_VERSION=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument."
      exit 1
      ;;
  esac
done

if [ -z $CHAINCODE_LANG ]; then
  CHAINCODE_LANG="golang"
fi

if [ -z $CHAINCODE_CONSTRUCTOR ]; then
  CHAINCODE_CONSTRUCTOR="[]"
fi

if [[ ! -z $CHAINCODEID && ! -z $CHAINCODE_VERSION ]]; then
  constructor="{\"Args\":$CHAINCODE_CONSTRUCTOR}"
  echo "INSTANTIATING chaincode $CHAINCODEID in $CHAINCODE_LANG on $ORDERER with version $CHAINCODE_VERSION in $CHANNELNAME with $ORDERER_CA"
  echo "with constructor $constructor"
  echo
  peer chaincode instantiate -o $ORDERER -C $CHANNELNAME -n $CHAINCODEID -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c $constructor  -P "OR ('Org1MSP.member', 'Org2MSP.member')" --tls --cafile $ORDERER_CA
else
  usage
fi
`

func generateInstantiatieChaincodeScript(spec NetworkSpec) error {

	scriptBody := []byte(instantiatieChaincodeScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "instantiate-chaincode.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}

	return nil
}

const upgradeChaincodeScriptTxt = `#!/bin/bash

function usage(){ 
  echo "Usage: $0 <arguments>"
  echo "Usage: $0 <flags>"
  echo "Mandatory:"
  echo "   -c <cc id>         A unique string identifier"
  echo "   -v <cc version>    A numeric number"
  echo "   -n <channelname>   A string identifying a channel"
  echo "Optional:"
  echo "   -a <cc argment>   <cc argument> must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
  echo "   -l <language>     Option to deloy in language other Go, currenly only node"
}

if [ "$#" -eq "0" ]; then  
  usage
  exit
fi

while getopts "a:c:l:n:v:" opt; do
  case $opt in
    a)
      CHAINCODE_CONSTRUCTOR=$OPTARG
      ;;
    c)
      CHAINCODEID=$OPTARG
      ;;
    l)
      CHAINCODE_LANG=$OPTARG
      ;;
    n)
      CHANNELNAME=$OPTAG
      ;;
    v)
      CHAINCODE_VERSION=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

if [ -z $CHAINCODE_CONSTRUCTOR ]; then
  CHAINCODE_CONSTRUCTOR="[]"
fi

if [ -z $CHAINCODE_LANG ]; then
  CHAINCODE_LANG="golang"
fi

if [[ ! -z $CHAINCODE_VERSION && ! -z $CHAINCODEID ]]; then

  echo "UPGRADING chaincode $CHAINCODEID to version $CHAINCODE_VERSION"
  echo "in $CHANNELNAME"
  echo "with constructor $CHAINCODE_CONSTRUCTOR"
  constructor="{\"Args\":$CHAINCODE_CONSTRUCTOR}"

  peer chaincode upgrade -o $ORDERER -C $CHANNELNAME -n $CHAINCODEID -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c $constructor --tls --cafile $ORDERER_CA
else
  usage
fi
`

func generateUpgradeChaincodeScript(spec NetworkSpec) error {

	scriptBody := []byte(upgradeChaincodeScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "upgrade-chaincode.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}

	return nil
}

const invokeScriptTxt = `#!/bin/bash

function usage(){ 
  echo "Usage: $0 <flags>"
  echo "Mandatory:"
  echo "   -c <cc id>        A unique string identifier"
  echo "   -n <channelname>  A string identifying a channel
  echo "Optional:"
  echo "   -a <cc constructor>   Must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
}

if [ "$#" -eq "0" ]; then  
  usage
  exit 1
fi

while getopts "a:c:n:" opt; do
  case $opt in
    a)
      CHAINCODE_CONSTRUCTOR=$OPTARG
      ;;
    c)
      CHAINCODEID=$OPTARG
      ;;
    n)
      CHANNELNAME=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

if [ -z $CHAINCODE_CONSTRUCTOR ]; then
  CHAINCODE_CONSTRUCTOR="[]"
fi

if [ ! -z $CHAINCODEID ]; then
  echo "INVOKING chaincode $CHAINCODEID in $CHANNELNAME on the $ORDERER"
  constructor="{\"Args\":$CHAINCODE_CONSTRUCTOR}"
  echo "with constructor $constructor"

  peer chaincode invoke -o $ORDERER -C $CHANNELNAME -n $CHAINCODEID -c $constructor
else
  usage
fi 
`

func generateInvokeScript(spec NetworkSpec) error {

	scriptBody := []byte(invokeScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "invoke.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}

	return nil
}

const queryScriptTxt = `#!/bin/bash

function usage(){ 
  echo "Usage: $0 <flags>"
  echo "Mandatory:"
  echo "   -c <cc id>        A unique string identifier"
  echo "   -n <channelname>  A string identifying a channel
  echo "Optional:"
  echo "   -a <cc constructor>   Must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
}

if [ "$#" -eq "0" ]; then  
  usage
  exit 1
fi

while getopts "a:c:n:" opt; do
  case $opt in
    a)
      CHAINCODE_CONSTRUCTOR=$OPTARG
      ;;
    c)
      CHAINCODEID=$OPTARG
      ;;
    n)
      CHANNELNAME=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

if [ -z $CHAINCODE_CONSTRUCTOR ]; then
   CHAINCODE_CONSTRUCTOR="[]"
fi

if [ ! -z $CHAINCODEID ]; then
  echo "Querying chaincode $CHAINCODEID in $CHANNELNAME on the $ORDERER"
  constructor="{\"Args\":$CHAINCODE_CONSTRUCTOR}"
  echo "with constructor $constructor"

  peer chaincode query -C $CHANNELNAME -n $CHAINCODEID -c $constructor
else
  usage
fi 
`

func generateQueryScript(spec NetworkSpec) error {

	scriptBody := []byte(queryScriptTxt)
	bashScript := filepath.Join(spec.ScriptPath, "query.sh")
	if err := ioutil.WriteFile(bashScript, scriptBody, 0777); err != nil {
		return err
	}

	return nil
}

// GenerateScripts produces scripts
func GenerateScripts(spec NetworkSpec) error {
	if err := generateCreateChannelScript(spec); err != nil {
		return err
	}
	if err := generateJoinChannelScript(spec); err != nil {
		return err
	}
	if err := generateInstallChaincodeScript(spec); err != nil {
		return err
	}
	if err := generateInstantiatieChaincodeScript(spec); err != nil {
		return err
	}
	if err := generateInvokeScript(spec); err != nil {
		return err
	}
	if err := generateQueryScript(spec); err != nil {
		return err
	}
	return nil
}
