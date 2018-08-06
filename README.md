# Introduction

The artefacts, or Fabric development kit, found in this Git repo is intended to help developers, using macOS or Linux-based platforms:

* learn what is involved in orchestrating and instantiating a Fabric network;
* debug chaincode;
* Customise a locally instantiable Fabric network to minic a production version.

The current version of this Fabric development kit does not yet support the Windows platform or to customise a Farbic network. This may be incorporated in future version.

# Pre-requisites

1. Install [Go tools](http://golang.org/dl).

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

# Two-organisations consortium

This configuration of Fabric network provides developers with the opportunity to quickly debug a chaincode. A default two-organisations network is provided where developers could subject their chaincode to interaction between two organisations.

Please refer to [this for a detail description](./docs/two-orgs.md)

# Acknowledgement

Aladdin Blockchain Technologies Ltd for sponsoring the effort to create this Fabric Development Kit 

# Disclaimer

Unless otherwise specified, the artefacts in this repository are distributed under Apache 2 license. In particular, the chaincodes are provided on "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

