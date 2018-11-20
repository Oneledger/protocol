/*
	Copyright 2017-2018 OneLedger
*/
package vm

type OLVMService struct {
	Protocol string
	//Port     int // TODO: Should be a full address (even if we only need port)
	Address string
}

type OLVMClient struct {
	Protocol    string
	ServicePath string // TODO: Should be called Address
}

type Container int

type Args struct {
	From       string
	Address    string
	CallString string
	Value      int
}

type Reply struct {
	Out string
	Ret string // TODO: Should be a real name
}
