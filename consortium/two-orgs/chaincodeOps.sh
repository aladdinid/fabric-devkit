#!/bin/bash

ARGS_NUMBER="$#"
COMMAND="$1"
ARG_2="$2"

usage_message="Useage: $0 install | instantiate | update | invoke | query | log <version number>"

if [ $ARGS_NUMBER -gt 1 ] && [ $ARGS_NUMBER -lt 2 ]; then
    echo $usage_message
    exit
fi

# If you plan to use node chaincode, replace the -l argument with `node` and change -p argument to an appropriate package name.
function installChaincode(){
    docker exec cli.peer0.org1.fabric.network /bin/bash -c '${PWD}/scripts/install-chaincode.sh -c mycc -v 1.0 -l golang -p minimalcc'
    docker exec cli.peer0.org2.fabric.network /bin/bash -c '${PWD}/scripts/install-chaincode.sh -c mycc -v 1.0 -l golang -p minimalcc'
}

# If you plan to use Go chaincode, replace the -l argument with `golang`
function instantiateChaincode(){
    docker exec cli.peer0.org1.fabric.network /bin/bash -c '${PWD}/scripts/instantiate-chaincode.sh -c mycc -v 1.0 -l golang -a [\"Init\",\"Paul\",\"10\",\"John\",\"20\"]'
}

# If you plan to use Go chaincode, replace the -l argument with `golang`
function updateChaincode(){
    docker exec cli.peer0.org1.fabric.network /bin/bash -c '${PWD}/scripts/upgrade-chaincode.sh -c mycc -v 1.0 -l golang'
}

function invokeChaincode(){
    docker exec cli.peer0.org1.fabric.network /bin/bash -c '${PWD}/scripts/invoke.sh -c mycc -a [\"move\",\"Paul\",\"John\",\"1\"]'
}

function queryChaincode(){
    docker exec cli.peer0.org1.fabric.network /bin/bash -c '${PWD}/scripts/query.sh -c mycc -a [\"query\",\"Paul\"]'
    docker exec cli.peer0.org1.fabric.network /bin/bash -c '${PWD}/scripts/query.sh -c mycc -a [\"query\",\"John\"]'
}

function logChaincode(){
    version=$1
    if [ -z $version ]; then
        echo $usage_message
        exit 1
    fi 
    echo "Log for peer0 org1"
    echo
    docker logs dev-peer0.org1.fabric.network.cn-mycc-$version
    echo "================================================="
    echo
    echo "Log for peer0 org2"
    echo
    docker logs dev-peer0.org2.fabric.network-mycc-$version
    echo "================================================="
    
}

case $COMMAND in
    "install")
        installChaincode
        ;;
    "instantiate")
        instantiateChaincode
        ;;
    "update")
        updateChaincode
        ;;
    "invoke")
        invokeChaincode
        ;;
    "query")
        queryChaincode
        ;;
    "log")
        logChaincode $ARG_2
        ;;
    *)
        echo $usage_message
        exit 1
esac
