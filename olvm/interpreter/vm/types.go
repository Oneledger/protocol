/*
	Copyright 2017-2018 OneLedger
*/
package vm

// Static information about the service parameters
type OLVMService struct {
	Protocol string
	//Port     int // TODO: Should be a full address (even if we only need port)
	Address string
}

// Static information about the client parameters
type OLVMClient struct {
	Protocol    string
	ServicePath string // TODO: Should be called Address
}

// TODO Still used?
type Container int
