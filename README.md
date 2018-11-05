# Aladdin Fabric Development Kit

Aladdin Fabric development kit, codename `Maejor`, was created to support developers building [Aladdin Blockchain Technologies Ltd](https://aladdinid.com/) a Fabric based [Proof-of-Concept (PoC) medical audit trail](https://www.youtube.com/watch?v=vJmhwymh-eU). Following from the PoC, Aladdin has decided to open source the developer kit.

The development kit was intended to help developers:

* learn, through experimentation, the system architecture of Hyperledger Fabric (Fabric);
* define, instantiate and configure a Fabric network to support development effort;
* verify that apps developed to interact with the Fabric network meets expectation.

The development kit has three components:

* Example Go chaincodes
* Reference implementations
* Command line interface (CLI) known as `maejor`

### Example Go chaincodes

## Acknowledgement

The maintainer(s) of this project is grateful to Aladdin for sponsoring the effort to create this developer kit.

## Release notes

### v0.1

Features:

* Implementation named `twoorgs`.
* Network components based on Fabric 1.1.0 containers and docker compose
* Bash based scripts to support kit orchestration - i.e. macOS and Ubuntu based

Status: 

* Released
* No further update

### v0.2

Features:

* Add a new reference implementation named `multichain`
* `meajor` - generating a skeleton Fabric network specifications

Status:

* Under development.

## Content

| Item | Description |
| --- | --- |
| `chaincodes` | Example Go chaincodes |
| `guides` | How-to documentation |
| `maejor` | Source codes for `maejor` application |
| `reference` | A collection of smoke testing implementation of a fabric network setup |

### Chaincodes

| Item | Description |
| --- | --- |
| `minimalcc` | This is a simple chaincode illustrating the transfer to some numeric value between two parties. |
| `one` | This is a chaincode intended to be used in conjuction with `multichain` reference implementation. |
| `two` | This is a chaincode intended to be used in conjunction with `multichain` reference implementation. |

### Guides

| Item | Description |
| --- | --- |
| `two-orgs.md` | User and contextual guide for the reference/twoorgs implementation |

### Maejor cli application

Please refer to [README](./maejor/README.md) for implementation information.

### Reference

| Item | Description |
| --- | --- |
| `multichain` | This is an example of an implementation involving two organisational nodes with two channels. It also based on the very latest Fabric implementation. |
| `twoorgs` | This is a full end-to-end example based on two organisational nodes, one channel, Fabric Node SDK, Hyperledger Explorer and Fabric 1.1.0. |

## Contributions and feedback

The maintainer(s) of this project welcomes feedback and contribution from anyone. However, to manage the development lifecycle, Aladdin Blockchain Technologies Ltd designated maintainer(s) shall retain sole descretion in deciding any features to be incorporated or removed from this repository. 

## Disclaimer

Unless otherwise specified, the artefacts in this repository are distributed under Apache 2 license.

All artefacts found here are provided on "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.