/*

   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/tendermint/tendermint/abci/types"
	"os"
)

func (app *App) Query(req RequestQuery) ResponseQuery {
	app.logger.Debug("ABCI: Query", "req", req, "path", req.Path, "data", req.Data)

	routerReq := NewRequestFromData(req.Path, req.Data)
	routerReq.Ctx = app.Context

	resp := &Response{}
	app.r.Handle(*routerReq, resp)

	// TODO proper response error handling
	result := ResponseQuery{
		Code:  types.CodeTypeOK,
		Index: 0, // TODO: What is this for?

		Log:  "Log Information",
		Info: "Info Information",

		Key:   action.Message("result"),
		Value: resp.Data,

		Proof:  nil,
		Height: int64(app.Context.balances.Version),
	}

	app.logger.Debug("ABCI: Query Result", "result", result)
	return result
}

func NewABCIRouter() Router {
	r := NewRouter("abci")

	fatalIfError(r.AddHandler("/account/list", GetAccount))
	fatalIfError(r.AddHandler("/account/add", GetAccount))
	fatalIfError(r.AddHandler("/account/delete", GetAccount))

	fatalIfError(r.AddHandler("/query/balance", GetBalance))

}

/*
		Handlers start here
 */
func GetBalance(req Request, resp *Response) {
	req.Parse()

	key := req.GetBytes("key")
	if len(key) == 0 {
		resp.Error("parameter key missing")
		return
	}

	accBalance := req.Ctx.balances.Get(key, true)
	resp.SetDataObj(accBalance)
}

// GetAccount by the name
func GetAccount(req Request, resp *Response) {
	req.Parse()

	name := req.GetString("name")
	if name == "" {
		resp.Error("parameter name missing")
		return
	}

	// TODO get account by name

}


/*
		utils
 */
func fatalIfError(err error) {
	if err != nil {
		// log.Fatal(err)
		os.Exit(1)
	}
}
