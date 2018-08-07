# Introduction

The artefacts here are hardcoded to instantiate a two-organisations fabric network with in-built RESTful APIs and deploy chaincode packaged in `minimalcc`.

---
**Disclaimer**

This version of the Fabric network is intended primarily to illustrate and educate users. It could also be used to support basic debugging of chaincode.

It is **not** to be replicated in a mission critical or production environment. If you wish to do so, you do so at your own risk.
---

# Content

| Item | Description |
| --- | --- |
| api/ | This folder contains nodeJS Express code for generating a RESTFul API to enable developers to invoke chaincode via a RESTful client such as Postman. |
| assets/ | Cryptographic and Fabric channel configuration assets. These are generated from executing the `fabricOps.sh`. |
| scripts/ | A series of low-level Bash shell scripts for configuring the Fabric network, installing and instantiating chaincodes. Note: These scripts are intended to be executed from within the `cli` containers based on `fabric-tools` images provided by Hyperledger. |
| chaincodeOps.sh | A Bash shell script to aggregate the low-level scripts to support installation, instantiation and invocation of chaincode. Note: this script is hard-coded to work with chaincode named `mycc` and `verion 1.0` only. |
| channelOps.sh | A Bash shell script to aggregate the low-level scripts to create a channel name `mychannel`. |
| `configtx.yaml` | A script used to generate artefacts for configuration a Fabric channel. |
| `crypto-config.yaml` | A script used to generate cryptographic materials to support the creation of Member Service Provider. These are hardcoded to support the two-organisations configuration defined in `network-config.yaml` |
| fabricOps.sh | A Bash script to support the operations to prepare artefacts for instantiation of a Fabric network and deployment of chaincodes. |
| generate-chainconfig.sh | A component of `fabricOps.sh` to support the creation of artefacts for the channel configuration. |
| generate-crypto.sh | A component of `fabricOps.sh` to support the creation of cryptographic artefacts. |
| `network-config.yaml` | The script for configuring docker containers used to instantiate Fabric network and support the creation of supporting artefacts. |

# How to use two-organisations consortium

## Context
The consortium is hard coded to demonstrate payment between two hypothetical entities (`Paul` and `John`). 

At the instantiation of the chaincode, the `Paul` is initialised with `10` units of unspecified value, and `John` is initialised with `20` units.

A transaction to initiate payment from `John` to `Paul`.

## Steps

1. Open a Terminal.
2. Navigate ( `cd` ) into `./consortium/twoorgs`.
3. Run the command `./fabricOps.sh init` - NOTE: on Ubuntu/Linux you may need to run the command `sudo ./fabricOps.sh init`.
4. Verify that you see operations to download images operations to create genesis.block and that there is no error. Also verify that these folders are also created: `./assets/channel-artefacts` and `./assets/crypto-config`.
5. Run the command `./fabricOps.sh start-network`
6. Run the command `./fabricOps.sh status` and verify if items (a) are all Up
7. Run the command `./fabricOps.sh configure-network`
8. Verify that there is no error.
9. Use your RESTful Client (e.g. postman) and execute a POST `http://localhost:8081/invoke` with this body `{ "fcn":"pay", "args":["Paul","1","John"] }` -- i.e. read pay `Paul` amount `1` from `John` .
10. Verify that you get a successful response.
11. Run the command `docker logs dev-peer0.org1.fabric.network-mycc-1.0` and verify that you see Item (B) -- i.e. this proves that the transactions is logged in `Org1`
12. Run the command `docker logs dev-peer0.org2.fabric.network-mycc-1.0` and verify that you see Item (C) -- i.e. this proves that the transactions is logged in `Org2`

---
**Note**

The Fabric network uses a piece of technology known as Docker container for instantiation. If you encounter any error as you step through the process, you should try to reset the network by manipulating docker containers. Use either one of the following commands:

* `fabricOps.sh clean` - this is a gentlier way of resetting the network. It will only attempt to reset aspects of the artefacts related to the network and leave any containers not related the Fabric network untouch. For example, if you have a none Fabric container that you are using for other purposes.
* `fabricOps.sh cleanall` - this is a catch all way of reseting your network. It will clean all Fabric and non-Fabric containers. Sometimes, you might need to do this because existing configuration may be sticky and you need to have a total clean setting.
* Using a combination of docker commands to selectively remove containers that underpin the operations of the Fabric network, e.g. `docker ps` and `docker rm`, etc. Please refer to documentation from `Docker`.
---

Item (A):
```
api.org1.fabric.network
api.org2.fabric.network
cli.peer0.org1.fabric.network
cli.peer0.org2.fabric.network
peer0.org1.fabric.network
peer0.org2.fabric.network
peer0.db.org1.fabric.network
peer0.db.org2.fabric.network
ca.org1.fabric.network
ca.org2.fabric.network
orderer.fabric.network
```

Item (B):
```
<time stamp> [minimalcc] Info -> INFO 001 Hello Init
<time stamp> [minimalcc] Infof -> INFO 002 Name1: Paul Amount1: 10
<time stamp> [minimalcc] Infof -> INFO 003 Name2: John Amount2: 20
<time stamp> [minimalcc] Info -> INFO 004 Hello Invoke
<time stamp> [minimalcc] Infof -> INFO 005 Pay: Paul amount: 1 from: John
<time stamp> [minimalcc] Infof -> INFO 006 Before payment - payeeCurrentState: 10
<time stamp> [minimalcc] Infof -> INFO 007 Before payment - payerCurrentState: 20
<time stamp> [minimalcc] Infof -> INFO 008 After payment - payee state: 11 payer state: 19
```

Item (C):
```
<time stamp> [minimalcc] Info -> INFO 004 Hello Invoke
<time stamp> [minimalcc] Infof -> INFO 005 Pay: Paul amount: 1 from: John
<time stamp> [minimalcc] Infof -> INFO 006 Before payment - payeeCurrentState: 10
<time stamp> [minimalcc] Infof -> INFO 007 Before payment - payerCurrentState: 20
<time stamp> [minimalcc] Infof -> INFO 008 After payment - payee state: 11 payer state: 19
```
