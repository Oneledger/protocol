package ons

import (
	"math/big"
	"os"

	"github.com/tendermint/tendermint/libs/db"

	"github.com/tendermint/tendermint/abci/types"

	"github.com/Oneledger/protocol/data/ons"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// global setup
func setup() []string {
	testDBs := []string{"test_dbpath"}
	return testDBs
}

// remove test db dir after
func teardown(dbPaths []string) {
	for _, v := range dbPaths {
		err := os.RemoveAll(v)
		if err != nil {
			errors.New("Remove test db file error")
		}
	}
}

func generateKeyPair() (crypto.Address, crypto.PubKey, ed25519.PrivKeyEd25519) {

	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr, pubkey, prikey
}

func assemblyCtxData(currencyName string, currencyDecimal int64, setBalanceStore bool, setLogger bool, setInitCoin bool,
	setCoinAddr crypto.Address, domainStore bool, setDomaintoStore bool, setDomainPrice int64, onsaleFlag bool, headerHeight int64) *action.Context {

	header := &types.Header{
		Height: headerHeight,
	}
	ctx := &action.Context{
		Header: header,
	}

	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	// balance store
	var store *balance.Store
	if setBalanceStore {

		store = balance.NewStore("tb", cs)
		ctx.Balances = store
	}
	// logger
	if setLogger {
		ctx.Logger = new(log.Logger)
	}
	// domainStore
	if domainStore {
		ds := ons.NewDomainStore("td", cs)
		ctx.Domains = ds

		if setDomaintoStore {
			currency := &balance.Currency{
				Name:    "OLT",
				Chain:   chain.Type(1),
				Decimal: int64(5),
			}
			coin := &balance.Coin{
				Currency: *currency,
				Amount:   balance.NewAmount(setDomainPrice),
			}
			d := &ons.Domain{
				Name:           "test_domain",
				OwnerAddress:   setCoinAddr.Bytes(),
				AccountAddress: setCoinAddr.Bytes(),
				OnSaleFlag:     onsaleFlag,
				ActiveFlag:     true,
				SalePrice:      *coin,
			}
			err := ctx.Domains.Set(d)
			ctx.Domains.State.Commit()
			if err != nil {
				errors.New("error happens when pre-setup a test domain")
			}
		}
	}
	// currencyList
	if currencyName != "" {
		// register new token
		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    currencyName,
			Chain:   chain.Type(1),
			Decimal: currencyDecimal,
		}
		err := currencyList.Register(currency)
		if err != nil {
			errors.New("register new token error")
		}
		ctx.Currencies = currencyList

		// set coin for account
		if setInitCoin {
			coin := balance.Coin{
				Currency: currency,
				Amount:   balance.NewAmount(1000000000000000000),
			}
			ba := balance.NewBalance()
			ba = ba.AddCoin(coin)
			err = store.Set(setCoinAddr.Bytes(), *ba)
			if err != nil {
				errors.New("setup testing token balance error")
			}
			store.State.Commit()
			ctx.Balances = store
		}
	}
	return ctx
}

func assemblyCreateDomainTx(owner crypto.Address, amount *action.Amount) *DomainCreate {

	createDoamin := &DomainCreate{
		Owner:   owner.Bytes(),
		Account: owner.Bytes(),
		Name:    "test_domain",
		Price:   *amount,
	}
	return createDoamin
}

func assemblyPurchaseDomainTx(owner crypto.Address, amount *action.Amount) *DomainPurchase {

	purchaseDomain := &DomainPurchase{
		Name:     "test_domain",
		Buyer:    owner.Bytes(),
		Account:  owner.Bytes(),
		Offering: *amount,
	}
	return purchaseDomain
}

func assemblySaleDomainTx(owner crypto.Address, amount *action.Amount, cancelSale bool) *DomainSale {

	saleDomain := &DomainSale{
		DomainName:   "test_domain",
		OwnerAddress: owner.Bytes(),
		Price:        *amount,
		CancelSale:   cancelSale,
	}
	return saleDomain
}
func assemblySendDomainTx(owner crypto.Address, amount *action.Amount) *DomainSend {

	sendDomain := &DomainSend{
		From:       owner.Bytes(),
		DomainName: "test_domain",
		Amount:     *amount,
	}
	return sendDomain
}

func assemblyUpdateDomainTx(owner crypto.Address, active bool) *DomainUpdate {

	updateDomain := &DomainUpdate{
		Owner:   owner.Bytes(),
		Account: owner.Bytes(),
		Name:    "test_domain",
		Active:  active,
	}
	return updateDomain
}

func assemblyTxData(replaceAddre bool, initTokenAmount string, txType string, failTypeCast bool, cancelSale bool, activeFlag bool) (action.SignedTx, crypto.Address) {

	owner, ownerPubkey, ownerPrikey := generateKeyPair()
	if replaceAddre {
		owner = crypto.Address("")
	}
	amt, _ := balance.NewAmountFromString(initTokenAmount, 10)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}
	fee := action.Fee{
		Price: action.Amount{Currency: "OLT", Value: balance.Amount{Int: *big.NewInt(0).Div(&amt.Int, big.NewInt(1000))}},
		Gas:   int64(10),
	}

	var tx Ons

	if failTypeCast {
		tx = &DomainCreate{
			Owner: owner.Bytes(),
		}
	} else {
		switch txType {
		case "create":
			tx = assemblyCreateDomainTx(owner, amount)
		case "purchase":
			tx = assemblyPurchaseDomainTx(owner, amount)
		case "sale":
			tx = assemblySaleDomainTx(owner, amount, cancelSale)
		case "send":
			tx = assemblySendDomainTx(owner, amount)
		case "update":
			tx = assemblyUpdateDomainTx(owner, activeFlag)
		}
	}
	txData, err := tx.Marshal()
	if err != nil {
		errors.New(err.Error())
	}
	rawtx := action.RawTx{
		Type: tx.Type(),
		Data: txData,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := ownerPrikey.Sign(rawtx.RawBytes())

	signed := action.SignedTx{
		RawTx: rawtx,
		Signatures: []action.Signature{
			action.Signature{
				Signer: keys.PublicKey{keys.ED25519, ownerPubkey.Bytes()[5:]},
				Signed: signature,
			}},
	}
	return signed, owner
}

func assemblyTxDataSameKp(replaceAddre bool, initTokenAmount string, txType string, failTypeCast bool, owner crypto.Address, ownerPubkey crypto.PubKey, ownerPrikey ed25519.PrivKeyEd25519) action.SignedTx {

	if replaceAddre {
		owner = crypto.Address("")
	}
	amt, _ := balance.NewAmountFromString(initTokenAmount, 10)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}
	fee := action.Fee{
		Price: action.Amount{Currency: "OLT", Value: balance.Amount{Int: *big.NewInt(0).Div(&amt.Int, big.NewInt(1000))}},
		Gas:   int64(10),
	}

	var txData action.Msg

	if failTypeCast {
		txData = &DomainCreate{
			Owner: owner.Bytes(),
		}
	} else {
		switch txType {
		case "create":
			txData = assemblyCreateDomainTx(owner, amount)
		case "purchase":
			txData = assemblyPurchaseDomainTx(owner, amount)
		case "sale":
			txData = assemblySaleDomainTx(owner, amount, false)
		case "send":
			txData = assemblySendDomainTx(owner, amount)
		case "update":
			txData = assemblyUpdateDomainTx(owner, true)
		}
	}
	data, err := txData.Marshal()
	if err != nil {
		errors.New(err.Error())
	}
	tx := action.RawTx{
		Type: txData.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}

	signature, _ := ownerPrikey.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			action.Signature{
				Signer: keys.PublicKey{keys.ED25519, ownerPubkey.Bytes()[5:]},
				Signed: signature,
			}},
	}
	return signed
}
