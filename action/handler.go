package action

type Tx interface {
	//it should be able to validate a tx by itself, non-valid tx will be reject by the node directly
	//without going to node process
	Validate(ctx *Context, signedTx SignedTx) (bool, error)

	//check tx on the first node who receives it, if process check failed, the tx will not be broadcast
	//could store a version of checked value in storage, but not implemented now
	ProcessCheck(ctx *Context, tx RawTx) (bool, Response)

	//deliver tx on the chain by changing the storage values, which should only be committed at Application
	//commit stage
	ProcessDeliver(ctx *Context, tx RawTx) (bool, Response)

	//process the charge of fees
	ProcessFee(ctx *Context, signedTx SignedTx, start Gas, size Gas, gasUsed Gas) (bool, Response)
}

//used for unknow transaction in router or not registered ones
type unknownTx struct {
}

func (unknownTx) Validate(ctx *Context, signedTx SignedTx) (bool, error) {
	return false, ErrWrongTxType
}

func (unknownTx) ProcessCheck(ctx *Context, tx RawTx) (bool, Response) {
	return false, Response{}
}

func (unknownTx) ProcessDeliver(ctx *Context, tx RawTx) (bool, Response) {
	return false, Response{}
}

func (unknownTx) ProcessFee(ctx *Context, signedTx SignedTx, start Gas, size Gas, gasUsed Gas) (bool, Response) {
	return false, Response{}
}
