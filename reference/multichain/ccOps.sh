#!/bin/bash

CHANNEL_ONE_NAME=channelone
CHANNEL_TWO_NAME=channeltwo
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

FIRST_CHAINCODE_NAME="firstchaincode"
FIRST_CHAINCODE_SRC="github.com/chaincode/one"
SECOND_CHAINCODE_NAME="secondchaincode"
SECOND_CHAINCODE_SRC="github.com/chaincode/two"
CHAINCODE_VERSION="1.0"

# Install and instantiate first chaincode in org 1
docker exec cli peer chaincode install -n $FIRST_CHAINCODE_NAME -p $FIRST_CHAINCODE_SRC -v $CHAINCODE_VERSION
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_ONE_NAME -c '{"Args":[]}' -n $FIRST_CHAINCODE_NAME -v $CHAINCODE_VERSION -P "OR('Org1MSP.member', 'Org2MSP.member')"

# Install first chaincode in org 2
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode install -n $FIRST_CHAINCODE_NAME -p $FIRST_CHAINCODE_SRC -v $CHAINCODE_VERSION
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_ONE_NAME -c '{"Args":[]}' -n $FIRST_CHAINCODE_NAME -v $CHAINCODE_VERSION -P "OR('Org1MSP.member', 'Org2MSP.member')"

# Install and instantiate second chaincode
docker exec  -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode install -n $SECOND_CHAINCODE_NAME -p $SECOND_CHAINCODE_SRC -v $CHAINCODE_VERSION
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_TWO_NAME -c '{"Args":[]}' -n $SECOND_CHAINCODE_NAME -v $CHAINCODE_VERSION -P "OR('Org1MSP.member', 'Org2MSP.member')"