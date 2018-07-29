#!/bin/bash

. ./scripts/common.sh

usage_message="Useage: $0 start | status| clean"

ARGS_NUMBER="$#"
COMMAND="$1"

function verifyArg() {
    if [ $ARGS_NUMBER -ne 1 ]; then
        echo $usage_message
        exit 1;
    fi
}

function pullDockerImages(){

    for IMAGES in ca peer orderer ccenv tools; do
        docker pull hyperledger/fabric-$IMAGES:$FABRIC_VERSION
        docker tag hyperledger/fabric-$IMAGES:$FABRIC_VERSION hyperledger/fabric-$IMAGES
    done

    docker pull hyperledger/fabric-couchdb:$COUCHDB_VERSION
    docker tag hyperledger/fabric-couchdb:$COUCHDB_VERSION hyperledger/fabric-couchdb

}

function setupAssetGenerator(){
    docker-compose -f ./network-config.yaml up -d assets.generator.fabric.network 
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
   docker exec assets.generator.fabric.network /bin/bash -c '${PWD}/generate-crypto.sh'
}

function generateChannelArtifacts(){
    docker exec assets.generator.fabric.network /bin/bash -c '${PWD}/generate-chanconfig.sh'
}

function createSecretPrivKeys(){
    pushd ./assets/crypto-config/peerOrganizations/org1.fabric.network/ca
    PK=$(ls *_sk)
    mv $PK secret.key
    popd

    pushd ./assets/crypto-config/peerOrganizations/org2.fabric.network/ca
    PK=$(ls *_sk)
    mv $PK secret.key
    popd
}

function startNetwork(){
    docker-compose -f ./network-config.yaml up -d orderer.fabric.network
    
    docker-compose -f ./network-config.yaml up -d ca.org1.fabric.network
    docker-compose -f ./network-config.yaml up -d peer0.db.org1.fabric.network
    docker-compose -f ./network-config.yaml up -d peer0.org1.fabric.network
    docker-compose -f ./network-config.yaml up -d cli.peer0.org1.fabric.network
     
    docker-compose -f ./network-config.yaml up -d ca.org2.fabric.network
    docker-compose -f ./network-config.yaml up -d peer0.db.org2.fabric.network
    docker-compose -f ./network-config.yaml up -d peer0.org2.fabric.network
    docker-compose -f ./network-config.yaml up -d cli.peer0.org2.fabric.network
}

function cleanNetwork(){
    fabric_container=$(docker ps --format "{{.ID}}" --filter "name=fabric.network") 
    docker rm -f $fabric_container
    asset_generator=$(docker ps --format "{{.ID}}" --filter "name=asset_generator")
    docker rm -f $asset_generator
    docker rmi -f $(docker images | grep hyperledger | tr -s ' ' | cut -d ' ' -f 3)
}

function networkStatus(){
    docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Status}}" --filter "name=fabric.network"
}

verifyArg
case $COMMAND in
    "start")
        pullDockerImages
        setupAssetGenerator
        cleanAssets
        sleep 1s
        generateCryptoArtifacts
        generateChannelArtifacts
        createSecretPrivKeys
        sleep 1s
        startNetwork
        ;;
    "status")
        networkStatus
        ;;
    "clean")
        cleanNetwork
        cleanAssets
        ;;  				
    *)
        echo $usage_message
        exit 1
esac