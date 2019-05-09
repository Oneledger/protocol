package app

// ABCI methods
type infoServer func(RequestInfo) ResponseInfo
type optionSetter func(RequestSetOption) ResponseSetOption
type queryer func(RequestQuery) ResponseQuery
type txChecker func([]byte) ResponseCheckTx
type chainInitializer func(RequestInitChain) ResponseInitChain
type blockBeginner func(RequestBeginBlock) ResponseBeginBlock
type txDeliverer func([]byte) ResponseDeliverTx
type blockEnder func(RequestEndBlock) ResponseEndBlock
type commitor func() ResponseCommit

// abciController ensures that the implementing type can control an underlying ABCI app
type abciController interface {
	infoServer() infoServer
	optionSetter() optionSetter
	queryer() queryer
	txChecker() txChecker
	chainInitializer() chainInitializer
	blockBeginner() blockBeginner
	txDeliverer() txDeliverer
	blockEnder() blockEnder
	commitor() commitor
}

var _ ABCIApp = &ABCI{}

// ABCI is used as an input for creating a new node
type ABCI struct {
	infoServer       infoServer
	optionSetter     optionSetter
	queryer          queryer
	txChecker        txChecker
	chainInitializer chainInitializer
	blockBeginner    blockBeginner
	txDeliverer      txDeliverer
	blockEnder       blockEnder
	commitor         commitor
}

func (app *ABCI) Info(request RequestInfo) ResponseInfo {
	return app.infoServer(request)
}

func (app *ABCI) SetOption(request RequestSetOption) ResponseSetOption {
	return ResponseSetOption{}
}

func (app *ABCI) Query(request RequestQuery) ResponseQuery {
	return app.queryer(request)
}

func (app *ABCI) CheckTx(tx []byte) ResponseCheckTx {
	return app.txChecker(tx)
}

func (app *ABCI) InitChain(request RequestInitChain) ResponseInitChain {
	return app.chainInitializer(request)
}

func (app *ABCI) BeginBlock(request RequestBeginBlock) ResponseBeginBlock {
	return app.blockBeginner(request)
}

func (app *ABCI) DeliverTx(tx []byte) ResponseDeliverTx {
	return app.txDeliverer(tx)
}

func (app *ABCI) EndBlock(request RequestEndBlock) ResponseEndBlock {
	return app.blockEnder(request)
}

func (app *ABCI) Commit() ResponseCommit {
	return app.commitor()
}
