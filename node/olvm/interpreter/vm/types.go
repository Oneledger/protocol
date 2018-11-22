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

/*
// All of the input necessary to perform a computation on a transaction
type OLVMRequest struct {
	// TODO: Original Transaction
	// TODO: Last execution context
	// TODO: Scripts (if we can follow the includes and get all of them)
	// TODO: Data Handle (some way to call out for large data requests)

	From       string
	Address    string
	CallString string
	Value      int
}

// All of the output received from the computation
type OLVMResult struct {
	// TODO: Any subseqeunce transaction that needs to be broadcasted
	// TODO: Last execution context

	Out string
	Ret string // TODO: Should be a real name
}
*/
