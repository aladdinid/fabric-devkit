#!/bin/bash

CHANNEL_ONE_NAME=channelone
CHANNEL_ONE_PROFILE=ChannelOne
CHANNEL_TWO_NAME=channeltwo
CHANNEL_TWO_PROFILE=ChannelTwo
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

if [ -d ./channel-artifacts ]; then
    rm -rf ./channel-artifacts
fi

if [ -d ./crypto-config ]; then
    rm -rf ./crypto-config
fi

cryptogen generate --config=./crypto-config.yml --output="./crypto-config"

if [ ! -d ./channel-artifacts ]; then
    mkdir -p ./channel-artifacts
fi

configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block

configtxgen -profile ${CHANNEL_ONE_PROFILE} -outputCreateChannelTx ./channel-artifacts/${CHANNEL_ONE_NAME}.tx -channelID $CHANNEL_ONE_NAME
configtxgen -profile ${CHANNEL_TWO_PROFILE} -outputCreateChannelTx ./channel-artifacts/${CHANNEL_TWO_NAME}.tx -channelID $CHANNEL_TWO_NAME

# generate anchor peer for channelone transaction of org1 
configtxgen -profile ${CHANNEL_ONE_PROFILE} -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors_${CHANNEL_ONE_NAME}.tx -channelID $CHANNEL_ONE_NAME -asOrg Org1MSP

# generate anchor peer for channelone channel transaction of org2
configtxgen -profile ${CHANNEL_ONE_PROFILE} -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors_${CHANNEL_ONE_NAME}.tx -channelID $CHANNEL_ONE_NAME -asOrg Org2MSP

# generate anchor peer for channeltwo transaction of org2
configtxgen -profile ${CHANNEL_TWO_PROFILE} -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors_${CHANNEL_TWO_NAME}.tx -channelID $CHANNEL_TWO_NAME -asOrg Org2MSP