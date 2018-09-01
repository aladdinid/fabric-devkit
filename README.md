# Introduction

The artefacts, or Fabric development kit, found in this Git repo is intended to help developers:

* learn what is involved in orchestrating and instantiating a Hyperledger Fabric (Fabric) network;
* debug chaincode.

# Features

**Current version**: 0.1.1-alpha

## Version: 0.1.1-alpha

1. A simple end-to-end command line (CLI) solution based on a fixed two-organisations Fabric network.
2. Verified for used on macOS and Linux only
3. Compatible with Docker engine 18.06.0-ce and Docker compose 1.22.0
4. Hyperledger Blockchain Explorer integrated. This is based on Explorer version 0.3.1 AS-IS.
5. This is only a alpha primarily for preview purposes only.

## For future considerations

1. Feature to enable developer customise a locally instantiable Fabric network to minic a production version for testing.
2. Feature to enable developer deploy a development only Fabric network in a shared platform (i.e. cloud, etc.) to enable multiple developers to collaborate.

# Pre-requisites

1. Install [Go](http://golang.org/dl).

    * for macOS, we recommend installing via [homebrew](http://brew.sh/);

    * for other platforms please refer to [installation guide](https://golang.org/doc/install).

2. Set the environmental variable GOPATH to a reference a directory to host your Go source codes and binaries (i.e. Go workspace). For example,

    `export GOPATH=$HOME/go-projects`

3. Create a folder in `$GOPATH/src` and navigate to it. Under the folder clone this repository.

# Content

| Item | Description |
| --- | --- |
| chaincodes/ | This folder contains chaincodes |
| chaincodes/minimalcc | This folder is the container for a default version of a chaincode for illustration purposes only |
| consortium | This folder contains definitions of and scripts to orchestrate Fabric network |
| consortium/twoorgs | This is a default two organisations fabric network intended for illustration purposes only |

# Two-organisations consortium

This configuration of Fabric network provides developers with the opportunity to quickly debug a chaincode. A default two-organisations network is provided where developers could subject their chaincode to interaction between two organisations.

Please refer to [this for a detail description](./docs/two-orgs.md)

# Acknowledgement

Aladdin Blockchain Technologies Ltd for sponsoring the effort to create this Fabric Development Kit 

# Disclaimer

Unless otherwise specified, the artefacts in this repository are distributed under Apache 2 license. 

All artefacts found here are provided on "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

Aladdin Blockchain Technologies Ltd has the descretion in deciding any features to be incorporated or removed from this repository.

