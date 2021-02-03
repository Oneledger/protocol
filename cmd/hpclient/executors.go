package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/passport"
)

const TZero = "0001-01-01T00:00:00Z"

func CreateHospitalAdminRequest(node LoadTestNode, threadNo int, shared *SharedData, _ interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in Create Hospital Token : thread", threadNo, r)
		}
	}()
	newAcc, err := CreateAccount(&node, threadNo)
	if err != nil {
		return
	}
	request := client.CreateTokenRequest{
		User:             passport.UserID(String(5)),
		TokenTypeID:      passport.TokenTypeID(String(6)),
		TokenType:        passport.TokenHospital,
		TokenSubType:     passport.ScreenerGeneral,
		OwnerAddress:     newAcc.Address(),
		SuperUserAddress: node.superadmin,
		SuperUser:        passport.UserID(node.superadminName),
		CreationTime:     GetTime(),
	}
	reply := &client.CreateTokenReply{}
	CreateToken(request, reply, "tx.CreateHospitalAdminToken", "broadcast.TxAsync", &node, threadNo)
}

func CreateScreenerAdminRequest(node LoadTestNode, threadNo int, shared *SharedData, _ interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in Create Screener Token : thread", threadNo, r)
		}
	}()
	newAcc, err := CreateAccount(&node, threadNo)
	if err != nil {
		return
	}
	request := client.CreateTokenRequest{
		User:             passport.UserID(String(5)),
		TokenTypeID:      passport.TokenTypeID(String(6)),
		TokenType:        passport.TokenScreener,
		TokenSubType:     passport.ScreenerGeneral,
		OwnerAddress:     newAcc.Address(),
		SuperUserAddress: node.superadmin,
		SuperUser:        passport.UserID(node.superadminName),
		CreationTime:     GetTime(),
	}
	reply := &client.CreateTokenReply{}
	CreateToken(request, reply, "tx.CreateScreenerAdminToken", "broadcast.TxAsync", &node, threadNo)
}

func AddTestInfoInit(node LoadTestNode, threadNo int, shared *SharedData) (out interface{}, err error) {
	// setup a hospital admin
	hosAdmin, err := CreateAccount(&node, threadNo)
	if err != nil {
		return
	}
	request := client.CreateTokenRequest{
		User:             passport.UserID(String(64)),
		TokenTypeID:      passport.TokenTypeID(String(64)),
		TokenType:        passport.TokenHospital,
		TokenSubType:     passport.ScreenerGeneral,
		OwnerAddress:     hosAdmin.Address(),
		SuperUserAddress: node.superadmin,
		SuperUser:        passport.UserID(node.superadminName),
		CreationTime:     GetTime(),
	}
	reply := &client.CreateTokenReply{}
	token := CreateToken(request, reply, "tx.CreateHospitalAdminToken", "broadcast.TxCommit", &node, threadNo)
	out = *token
	return
}

func AddTestInfoRequest(node LoadTestNode, threadNo int, shared *SharedData, in interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in AddTestInfoRequest : thread", threadNo, r)
		}
	}()
	fullnode := node.clCtx.FullNodeClient()

	// create an account and random result
	acc, err := CreateAccount(&node, threadNo)
	if err != nil {
		return
	}
	result := passport.COVID19Pending

	// request upload
	token, _ := in.(passport.AuthToken)
	request := client.AddTestInfoRequest{
		TestID:       String(64),
		Person:       passport.UserID(String(64)),
		Test:         passport.TestCOVID19,
		SubTest:      passport.TestSubPCR,
		Manufacturer: "test manufacturer",
		Result:       result,

		TestOrg:     token.TokTypeId,
		TestedAt:    GetTime(),
		TestedBy:    String(64),
		AnalysisOrg: passport.TypeIDInvalid,
		AnalyzedAt:  TZero,
		AnalyzedBy:  "",

		Admin:        token.OwnerId,
		AdminAddress: token.OwnerAddress,
		UploadedAt:   GetTime(),
		Notes:        "some notes",
	}
	reply := &client.CreateTxReply{}
	err = fullnode.Client.Call("tx.AddTestInfo", request, reply)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error calling tx.AddTestInfo: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// sign Tx
	replySign := &client.SignRawTxResponse{}
	err = fullnode.Client.Call("owner.SignWithSecureAddress", client.SecureSignRawTxRequest{
		RawTx:    reply.RawTx,
		Address:  request.AdminAddress,
		Password: "1234",
		KeyPath:  node.keypath,
	}, replySign)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in Sign: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// broadcast Tx
	bresult := &client.BroadcastReply{}
	err = fullnode.Client.Call("broadcast.TxAsync", client.BroadcastRequest{
		RawTx:     reply.RawTx,
		Signature: replySign.Signature.Signed,
		PublicKey: replySign.Signature.Signer,
	}, bresult)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in BroadcastTxAsync:%s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
	if !bresult.OK {
		logger.Errorf("Thread : %d | Node : %s | BroadcastTxAsync Not ok:%s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// request update
	request.AnalysisOrg = token.TokTypeId
	request.AnalyzedAt = GetTime()
	request.AnalyzedBy = "Doctor Hasson"
	if rand.Intn(2) == 0 {
		request.Result = passport.COVID19Negative
	} else {
		request.Result = passport.COVID19Positive
	}
	reply = &client.CreateTxReply{}
	err = fullnode.Client.Call("tx.AddTestInfo", request, reply)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error calling tx.AddTestInfo(update): %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
	// sign update Tx
	replySign = &client.SignRawTxResponse{}
	err = fullnode.Client.Call("owner.SignWithSecureAddress", client.SecureSignRawTxRequest{
		RawTx:    reply.RawTx,
		Address:  request.AdminAddress,
		Password: "1234",
		KeyPath:  node.keypath,
	}, replySign)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in Sign(update): %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
	// broadcast update Tx
	bresult = &client.BroadcastReply{}
	err = fullnode.Client.Call("broadcast.TxAsync", client.BroadcastRequest{
		RawTx:     reply.RawTx,
		Signature: replySign.Signature.Signed,
		PublicKey: replySign.Signature.Signer,
	}, bresult)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in BroadcastTxAsync(update):%s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
	if !bresult.OK {
		logger.Errorf("Thread : %d | Node : %s | BroadcastTxAsync Not ok(update):%s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// add to list of tested person
	shared.Lock()
	defer shared.Unlock()
	shared.persons = append(shared.persons, request.Person)
	shared.personAddrs = append(shared.personAddrs, acc.Address())
	logger.Detailf("Thread : %d | Node : %s | Person Tested: %s", threadNo, node.cfg.Node.NodeName, request.Person)
}

func ReadTestInfoInit(node LoadTestNode, threadNo int, shared *SharedData) (out interface{}, err error) {
	// setup a screener admin
	newAcc, err := CreateAccount(&node, threadNo)
	if err != nil {
		return
	}
	request := client.CreateTokenRequest{
		User:             passport.UserID(String(64)),
		TokenTypeID:      passport.TokenTypeID(String(64)),
		TokenType:        passport.TokenScreener,
		TokenSubType:     passport.ScreenerGeneral,
		OwnerAddress:     newAcc.Address(),
		SuperUserAddress: node.superadmin,
		SuperUser:        passport.UserID(node.superadminName),
		CreationTime:     GetTime(),
	}
	reply := &client.CreateTokenReply{}
	token := CreateToken(request, reply, "tx.CreateScreenerAdminToken", "broadcast.TxCommit", &node, threadNo)
	out = *token
	return
}

func ReadTestInfoRequest(node LoadTestNode, threadNo int, shared *SharedData, in interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in Upload Test Info : thread", threadNo, r)
		}
	}()
	fullnode := node.clCtx.FullNodeClient()
	token, _ := in.(passport.AuthToken)

	// catch a random person to scan
	var person passport.UserID
	var addr keys.Address
	for {
		shared.Lock()
		if len(shared.persons) == 0 {
			shared.Unlock()
		} else {
			index := rand.Intn(len(shared.persons))
			person = shared.persons[index]
			addr = shared.personAddrs[index]
			shared.Unlock()

			// try read
			req := client.PSPTFilterTestInfoRequest{
				Org:    token.TokTypeId,
				Admin:  token.OwnerId,
				Test:   passport.TestCOVID19,
				Person: person,
			}
			rpy := &client.PSPTFilterTestInfoReply{}
			err := fullnode.Client.Call("query.PSPT_QueryTestInfoByID", req, rpy)
			if err != nil {
				logger.Errorf("Thread : %d | Node : %s | error calling query.PSPT_QueryTestInfoByID: %s", threadNo, node.cfg.Node.NodeName, err)
			}

			// until we find no-empty list
			if len(rpy.InfoList) > 0 {
				for _, info := range rpy.InfoList {
					logger.Detailf("TestInfo : %s", info.String())
				}
				break
			}
		}
		time.Sleep(time.Millisecond)
	}

	// request
	request := client.ReadTestInfoRequest{
		Org:          token.TokTypeId,
		Admin:        token.OwnerId,
		AdminAddress: token.OwnerAddress,
		Person:       person,
		Address:      addr,
		Test:         passport.TestCOVID19,
		ReadAt:       GetTime(),
	}
	reply := &client.CreateTxReply{}
	err := fullnode.Client.Call("tx.ReadTestInfo", request, reply)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error calling tx.ReadTestInfo: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// sign Tx
	replySign := &client.SignRawTxResponse{}
	err = fullnode.Client.Call("owner.SignWithSecureAddress", client.SecureSignRawTxRequest{
		RawTx:    reply.RawTx,
		Address:  request.AdminAddress,
		Password: "1234",
		KeyPath:  node.keypath,
	}, replySign)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | Address : %s | error in Sign: %s", threadNo, node.cfg.Node.NodeName, request.Address, err)
		return
	}

	// broadcast Tx
	bresult := &client.BroadcastReply{}
	err = fullnode.Client.Call("broadcast.TxAsync", client.BroadcastRequest{
		RawTx:     reply.RawTx,
		Signature: replySign.Signature.Signed,
		PublicKey: replySign.Signature.Signer,
	}, bresult)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in BroadcastTxAsync: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
	if !bresult.OK {
		logger.Errorf("Thread : %d | Node : %s | BroadcastTxAsync Not ok: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
}

func CreateAccount(node *LoadTestNode, threadNo int) (newAcc accounts.Account, err error) {
	accName := String(5)
	newAcc, err = accounts.GenerateNewAccount(chain.Type(1), accName)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | Error in Creating Account: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// add to wallet key store
	wallet, err := accounts.NewWalletKeyStore(filepath.Clean(node.keypath))
	if err != nil {
		return
	}
	if !wallet.Open(newAcc.Address(), "1234") {
		logger.Errorf("Thread : %d | Node : %s | Error in open wallet: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}
	err = wallet.Add(newAcc)
	wallet.Close()
	return
}

func CreateToken(req client.CreateTokenRequest, reply *client.CreateTokenReply,
	txService, broadcastService string, node *LoadTestNode, threadNo int) (token *passport.AuthToken) {
	// prepare Tx
	err := node.clCtx.FullNodeClient().Client.Call(txService, req, reply)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in %s: %s", threadNo, node.cfg.Node.NodeName, txService, err)
		return
	}
	replySign := &client.SignRawTxResponse{}
	err = node.clCtx.FullNodeClient().Client.Call("owner.SignWithSecureAddress", client.SecureSignRawTxRequest{
		RawTx:    reply.RawTx,
		Address:  node.superadmin,
		Password: "1234",
		KeyPath:  node.keypath,
	}, replySign)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s | error in Sign: %s", threadNo, node.cfg.Node.NodeName, err)
		return
	}

	// broadcast Tx
	bresult := &client.BroadcastReply{}
	err = node.clCtx.FullNodeClient().Client.Call(broadcastService, client.BroadcastRequest{
		RawTx:     reply.RawTx,
		Signature: replySign.Signature.Signed,
		PublicKey: replySign.Signature.Signer,
	}, bresult)
	if err != nil {
		logger.Errorf("Thread : %d | Node : %s |error in %s: %s", threadNo, node.cfg.Node.NodeName, broadcastService, err)
		return
	}
	if !bresult.OK {
		logger.Errorf("Thread : %d | Node : %s | %s Not ok: %s", threadNo, node.cfg.Node.NodeName, broadcastService, err)
		return
	}

	// return token
	token = &passport.AuthToken{
		TokType:      req.TokenType,
		TokSubType:   req.TokenSubType,
		TokTypeId:    req.TokenTypeID,
		OwnerId:      req.User,
		OwnerAddress: req.OwnerAddress,
	}
	return
}
