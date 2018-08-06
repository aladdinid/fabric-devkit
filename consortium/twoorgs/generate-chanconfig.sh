#!/bin/bash

. ./scripts/common.sh

if [ -d channel-artefacts ]; then
    rm -rf ./assets/channel-artefacts
fi

mkdir ./assets/channel-artefacts

configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./assets/channel-artefacts/genesis.block
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./assets/channel-artefacts/channel.tx -channelID $CHANNELNAME
