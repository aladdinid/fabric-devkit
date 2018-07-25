#!/bin/bash

. ./scripts/common.sh

function usage(){ 
    echo "Usage: $0 <flags>"
    echo "Mandatory:"
    echo "   -c <cc id>         A unique string identifier"
    echo "   -v <cc version>    A numeric number"
    echo "   -p <cc package>    A name of folder containing chaincodes"
    echo "Optional:"
    echo "   -a <cc constructor>  Must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
    echo "   -l <language>        Option to deloy in language other Go, currenly only node"
}

echo "the number of arguments are: $#"
if [ "$#" -eq "0" ]; then  
    usage
    exit
fi

while getopts "a:c:l:p:v:" opt; do
  case $opt in
    a)
      CHAINCODE_CONSTRUCTOR=$OPTARG
      ;;
    c)
      CHAINCODEID=$OPTARG
      ;;
    l)
      CHAINCODE_LANG=$OPTARG
      ;;
    p)
      CHAINCODE_PACKAGE=$OPTARG
      ;;
    v)
      CHAINCODE_VERSION=$OPTARG
      ;;
    \?)
      usage
      exit 1
      ;;
    :)
      usage
      exit 1
      ;;
  esac
done

if [ -z $CHAINCODE_LANG ]; then
  CHAINCODE_LANG="golang"
fi

if [ -z $CHAINCODE_CONSTRUCTOR ]; then
    CHAINCODE_CONSTRUCTOR="[]"
fi

if [[ ! -z $CHAINCODEID && ! -z $CHAINCODE_VERSION && ! -z $CHAINCODE_PACKAGE ]]; then

    if [ "$CHAINCODE_LANG" == "golang" ]; then 
      path_to_chaincode="github.com/hyperledger/fabric/chaincodes/$CHAINCODE_PACKAGE"
    else
      path_to_chaincode="$GOPATH/src/github.com/hyperledger/fabric/chaincodes/$CHAINCODE_PACKAGE"
    fi
    echo "INSTALLING chaincode $CHAINCODEID in $CHAINCODE_LANG version $CHAINCODE_VERSION in $path_to_chaincode"
    echo
    peer chaincode install -n $CHAINCODEID -v $CHAINCODE_VERSION -l $CHAINCODE_LANG -p $path_to_chaincode --tls --cafile $ORDERER_CA

    echo "INSTALLATION chaincode successfully"
else
    usage
fi
