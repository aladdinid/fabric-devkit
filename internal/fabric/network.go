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

package fabric

import (
	"os"
	"text/template"

	"github.com/aladdinid/fabric-devkit/internal/config"
)

const networkConfig = `
version: '2'

services:

  orderer.{{Domain}}:
    container_name: orderer.{{Domain}}
    image: hyperledger/fabric-orderer
	tty: true
    environment:
      - CONFIGTX_ORDERER_ORDERERTYPE=solo
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=7050
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_QUEUESIZE=1000
      - ORDERER_GENERAL_MAXWINDOWSIZE=1000
      - ORDERER_RAMLEDGER_HISTORY_SIZE=100
      - ORDERER_GENERAL_BATCHSIZE=10
      - ORDERER_GENERAL_BATCHTIMEOUT=10s
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/fabric/crypto-config/channel-artefacts/genesis.block
      - ORDERER_GENERAL_LOCALMSPID={{$orderer.MSPID}}
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/fabric/crypto-config/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/fabric/crypto-config/tls/server.crt
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/fabric/crypto-config/tls/server.key
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/fabric/crypto-config/tls/ca.crt, /var/hyperledger/fabric/crypto-config/peerOrganizations/org1.fabric.network/tls/ca.crt, /var/hyperledger/fabric/crypto-config/peerOrganizations/org2.fabric.network/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./assets/channel-artefacts/:/var/hyperledger/fabric/crypto-config/channel-artefacts/
      - ./assets/crypto-config/ordererOrganizations/fabric.network/orderers/orderer.fabric.network/:/var/hyperledger/fabric/crypto-config/
      - ./assets/crypto-config/peerOrganizations/org1.fabric.network/peers/peer0.org1.fabric.network/tls/ca.crt:/var/hyperledger/fabric/crypto-config/peerOrganizations/org1.fabric.network/tls/ca.crt
      - ./assets/crypto-config/peerOrganizations/org2.fabric.network/peers/peer0.org2.fabric.network/tls/ca.crt:/var/hyperledger/fabric/crypto-config/peerOrganizations/org2.fabric.network/tls/ca.crt
    ports:
      - 7050:7050
`

// GenerateNetwork produces docker compose network compose file
func GenerateNetwork(config config.NetworkSpec) {
	tpl := template.Must(template.New("Main").Parse(networkConfig))
	tpl.Execute(os.Stdout, config)
}
