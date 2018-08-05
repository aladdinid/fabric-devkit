package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("minimalcc")

// Transaction
func initialise(stub shim.ChaincodeStubInterface) pb.Response {

	args := stub.GetArgs()

	name1 := string(args[1])
	amount1 := args[2]

	logger.Infof("Name1: %s Amount1: %s", name1, string(amount1))

	err := stub.PutState(name1, amount1)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to store state %v", err))
	}

	name2 := string(args[3])
	amount2 := args[4]

	logger.Infof("Name2: %s Amount2: %s", name2, string(amount2))

	err = stub.PutState(name2, amount2)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to store state %v", err))
	}

	return shim.Success([]byte("Initialisation completed"))
}

// SimpleChaincode representing a class of chaincode
type SimpleChaincode struct{}

// Init to initiate the SimpleChaincode class
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Hello Init")
	fcn, _ := stub.GetFunctionAndParameters()
	if fcn == "init" {
		return initialise(stub)
	}

	return shim.Error("Fail to initialise state")
}

// Invoke a method specified in the SimpleChaincode class
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Hello Invoke")
	fcn, args := stub.GetFunctionAndParameters()
	logger.Infof("Function: %v Arguments: %v", fcn, args)
	return shim.Success([]byte("Invoke"))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Debugf("Error: %s", err)
	}
}
