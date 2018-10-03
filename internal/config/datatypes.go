package config

// OrdererSpec is an abstraction of a fabric orderer
type OrdererSpec struct {
	Name   string
	Domain string
	MSPID  string
}

// NetworkSpec is an abstraction of a fabric configuration
type NetworkSpec struct {
	Orderers []OrdererSpec
}
