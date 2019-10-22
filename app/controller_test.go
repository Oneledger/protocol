package app

import (
	"testing"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"

	"github.com/stretchr/testify/assert"
)

func setupForTx() action.SignedTx {
	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	send := transfer.Send{
		From:   keys.Address(addr),
		To:     keys.Address(addr),
		Amount: action.Amount{"OLT", *balance.NewAmount(10)},
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(10)},
		Gas:   int64(10),
	}
	signer := action.Signature{
		Signer: keys.PublicKey{keys.ED25519, pubkey.Bytes()},
	}
	data, _ := send.Marshal()

	tx := action.SignedTx{
		RawTx: action.RawTx{
			Data: data,
			Fee:  fee,
			Memo: "test_memo",
		},
		Signatures: []action.Signature{signer},
	}
	return tx
}

func TestApp_infoServer(t *testing.T) {
	app := setup(setupForGeneralInfo)
	req := RequestInfo{}

	infoServer := app.infoServer()
	is := infoServer(req)

	assert.NotEmpty(t, is)
}

// TODO can't pass setupState step
func TestApp_chainInitializer(t *testing.T) {
	app := setup(setupForStart)

	chainInitializer := app.chainInitializer()
	req := RequestInitChain{}
	ci := chainInitializer(req)

	assert.Empty(t, ci.Validators)
}

func TestApp_blockBeginner(t *testing.T) {
	app := setup(setupForStart)

	req := RequestBeginBlock{}
	blockBeginner := app.blockBeginner()
	bb := blockBeginner(req)

	assert.Equal(t, 0, bb.Size())
}

func TestApp_txChecker(t *testing.T) {
	app := setup(setupForStart)

	tx := setupForTx()
	rawTx, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	assert.Nil(t, err)

	txChecker := app.txChecker()
	tc := txChecker(rawTx)
	assert.False(t, tc.IsOK())
}

func TestApp_txDeliverer(t *testing.T) {
	app := setup(setupForStart)

	tx := setupForTx()
	rawTx, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	assert.Nil(t, err)

	txDeliverer := app.txDeliverer()
	td := txDeliverer(rawTx)

	assert.False(t, td.IsOK())
}

func TestApp_blockEnder(t *testing.T) {
	app := setup(setupForStart)

	req := RequestEndBlock{}
	blockEnder := app.blockEnder()
	be := blockEnder(req)

	assert.Empty(t, be.ValidatorUpdates)
	assert.Equal(t, 0, be.Size())
}

func TestApp_commitor(t *testing.T) {
	app := setup(setupForStart)

	commitor := app.commitor()
	c := commitor()
	if c.Size() == 0 {
		assert.Empty(t, c.Data)
	} else {
		assert.NotEmpty(t, c.Data)
	}
}
