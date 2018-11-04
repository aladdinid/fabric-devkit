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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const networkSpec = `{{- $domain := .Domain}} {{- $chaincodepath := .ChaincodePath}}
version: '2'

services:

  orderer.{{$domain}}:
    container_name: orderer.{{$domain}}
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
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./channel-artefacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ./crypto-config/ordererOrganizations/{{$domain}}/orderers/orderer.{{$domain}}/msp:/var/hyperledger/orderer/msp
      - ./crypto-config/ordererOrganizations/{{$domain}}/orderers/orderer.{{$domain}}/tls:/var/hyperledger/orderer/tls
    ports:
      - 7050:7050

{{range $index, $org := .OrganizationSpecs }}
  # {{$org.Name}}
  ca.{{$org.Name | ToLower}}.{{$domain}}:
    container_name: ca.{{$org.Name | ToLower}}.{{$domain}}
    image: hyperledger/fabric-ca
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca.org1.fabric.network
      - FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server/crypto-config/ca/ca.{{$org.Name | ToLower}}.{{$domain}}-cert.pem
      - FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server/crypto-config/ca/secret.key
      - FABRIC_CA_SERVER_TLS_ENABLED=true
      - FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server/crypto-config/ca/ca.{{$org.Name | ToLower}}.{{$domain}}-cert.pem
      - FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server/crypto-config/ca/secret.key
    command: sh -c 'fabric-ca-server start -b admin:adminpw -d'
    volumes:
      - ./crypto-config/peerOrganizations/{{$org.Name | ToLower}}.{{$domain}}/ca/:/etc/hyperledger/fabric-ca-server/crypto-config/ca/
    ports:
      - {{PortInc 7054 $index}}:7054

  {{$org.Anchor}}.db.{{$org.Name | ToLower}}.{{$domain}}:
    container_name: {{$org.Anchor}}.db.{{$org.Name | ToLower }}.{{$domain}}
    image: hyperledger/fabric-couchdb
    ports:
      - {{PortInc 5984 $index}}:5984         

  {{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}:
    container_name: {{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}
    image: hyperledger/fabric-peer
    tty: true
    environment:
      - CORE_PEER_ID={{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}
      - CORE_PEER_ADDRESS={{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}:7051
      - CORE_PEER_LOCALMSPID={{$org.Name}}MSP
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/msp
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/tls/ca.crt
      - CORE_PEER_ENDORSER_ENABLED=true
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT={{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}:7051
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      # This disables mutual auth for gossip
      - CORE_PEER_GOSSIP_SKIPHANDSHAKE=true
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_LOGGING_LEVEL=debug
      - CORE_NEXT=true
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=network_default
      - GOPATH=/opt/gopath
      # World state
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS={{$org.Anchor}}.db.{{$org.Name | ToLower}}.{{$domain}}:5984
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: peer node start 
             /bin/bash
    volumes:
      - /var/run/:/host/var/run/
      - ./crypto-config/peerOrganizations/{{$org.Name | ToLower}}.fabric.network/peers/{{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}/:/etc/hyperledger/fabric/crypto-config/
    ports:
      - {{PortInc 7051 $index}}:7051
      - {{PortInc 7053 $index}}:7053
    depends_on: 
      - orderer.{{$domain}}
      - {{$org.Anchor}}.db.{{$org.Name | ToLower}}.{{$domain}}
  
  cli.{{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}:
    container_name: cli.{{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}
    image: hyperledger/fabric-tools
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=debug
      - CORE_PEER_ID={{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}
      - CORE_PEER_ADDRESS={{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}:7051
      - CORE_PEER_LOCALMSPID={{$org.Name}}MSP
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/tls/ca.crt
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/users/Admin@{{$org.Name | ToLower}}.{{$domain}}/msp
      - ORDERER_CA=/etc/hyperledger/fabric/crypto-config/orderer/msp/tlscacerts/tlsca.{{$domain}}-cert.pem
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: /bin/bash
    volumes:
      - /var/run/:/host/var/run/
      - ./crypto-config/peerOrganizations/{{$org.Name | ToLower}}.{{$domain}}/peers/{{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}/:/etc/hyperledger/fabric/crypto-config/
      - ./crypto-config/peerOrganizations/{{$org.Name}}.{{$domain}}/users/:/etc/hyperledger/fabric/crypto-config/users/
      - ./crypto-config/ordererOrganizations/{{$domain}}/orderers/orderer.{{$domain}}/:/etc/hyperledger/fabric/crypto-config/orderer/
      - ./channel-artefacts/:/opt/gopath/src/github.com/hyperledger/fabric/channel-artefacts/
      - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/scripts
      - {{$chaincodepath}}:/opt/gopath/src/github.com/hyperledger/fabric/chaincodes/
    depends_on:
      - {{$org.Anchor}}.{{$org.Name | ToLower}}.{{$domain}}
{{end}}
`

// CreateNetworkSpec produces docker compose network compose file
func CreateNetworkSpec(spec NetworkSpec) error {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"PortInc": func(port, index int) string {
			p := port + index*10
			return strconv.Itoa(p)
		},
	}

	tpl := template.Must(template.New("Main").Funcs(funcMap).Parse(networkSpec))
	networkSpecYML := filepath.Join(spec.NetworkPath, "network-config.yaml")
	f, err := os.Create(networkSpecYML)
	if err != nil {
		return err
	}
	err = tpl.Execute(f, spec)
	if err != nil {
		return err
	}

	return nil
}
