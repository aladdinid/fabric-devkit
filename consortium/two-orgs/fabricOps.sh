#!/bin/bash

function pullDockerImages(){

    for IMAGES in tools; do
        docker pull hyperledger/fabric-$IMAGES:x86_64-1.1.0
        docker tag hyperledger/fabric-$IMAGES:x86_64-1.1.0 hyperledger/fabric-$IMAGES
    done

    docker-compose -f ./fabric-cli.yaml up -d

}

function cleanAssets(){
    if [ -d ./assets/channel-artefacts ]; then
        rm -R -f ./assets/channel-artefacts/
    fi

    if [ -d ./assets/crypto-config ]; then
        rm -R -f ./assets/crypto-config/
    fi
}

function generateCryptoArtifacts(){
   docker exec cli.fabric.network /bin/bash -c '${PWD}/generate-crypto.sh'
}

function generateChannelArtifacts(){
    docker exec cli.fabric.network /bin/bash -c '${PWD}/generate-chanconfig.sh'
}


pullDockerImages
cleanAssets
sleep 1s
generateCryptoArtifacts
generateChannelArtifacts