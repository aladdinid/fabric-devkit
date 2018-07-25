#!/bin/bash

. ./scripts/common.sh

function usage(){ 
    echo "Usage: $0 <flags>"
    echo "Mandatory:"
    echo "   -c <cc id>         A unique string identifier"
    echo "   -v <cc version>    A numeric number"
    echo "Optional:"
    echo "   -a <cc constructor>  Must be in the form [\"method\", \"method-arg-1\", \"method-arg-2\"]"
    echo "   -l <language>        Option to deloy in language other Go, currenly only node"
}

if [ "$#" -eq "0" ]; then  
    usage
    exit
fi

while getopts "a:c:l:v:" opt; do
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
    v)
      CHAINCODE_VERSION=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument."
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

if [[ ! -z $CHAINCODEID && ! -z $CHAINCODE_VERSION ]]; then

    constructor="{\"Args\":$CHAINCODE_CONSTRUCTOR}"
    echo "INSTANTIATING chaincode $CHAINCODEID in $CHAINCODE_LANG on $ORDERER with version $CHAINCODE_VERSION in $CHANNELNAME with $ORDERER_CA"
    echo "with constructor $constructor"
    echo
    peer chaincode instantiate -o $ORDERER -C $CHANNELNAME -n $CHAINCODEID -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c $constructor  -P "OR ('AladdinMSP.member')" --tls --cafile $ORDERER_CA
else
    usage
fi
