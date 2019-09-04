package ons

import (
	"testing"

	"github.com/Oneledger/protocol/data/ons"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/action"
)

var txTypeList map[string]testSuit

type testSuit struct {
	tx   action.Tx
	want bool
}

type txResult struct {
	createResult   bool
	purchaseResult bool
	saleResult     bool
	updateResult   bool
	sendResult     bool
}

func setupForTx(tr *txResult) {

	txTypeList = make(map[string]testSuit, 5)

	txTypeList["create"] = testSuit{&domainCreateTx{}, tr.createResult}
	txTypeList["purchase"] = testSuit{&domainPurchaseTx{}, tr.purchaseResult}
	txTypeList["sale"] = testSuit{&domainSaleTx{}, tr.saleResult}
	txTypeList["update"] = testSuit{&domainUpdateTx{}, tr.updateResult}
	txTypeList["send"] = testSuit{&domainSendTx{}, tr.sendResult}

}

func TestDomainTx_Validate(t *testing.T) {
	ctx := &action.Context{}
	t.Run("validate basic tx with invalid tx data, should return error for each tx type", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			tx, _ := assemblyTxData(true, "100", txType, false, false, true)
			ok, _ := txRefer.tx.Validate(ctx, tx)
			assert.Equal(t, txRefer.want, ok, "%s domain Validate fails ", txType)
		}
	})
	t.Run("cast tx with invalid tx data, should return error for each tx type", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			tx, _ := assemblyTxData(false, "100", txType, true, false, true)
			ok, _ := txRefer.tx.Validate(ctx, tx)
			assert.Equal(t, txRefer.want, ok, "%s domain Validate fails ", txType)
		}

	})
	t.Run("invalid amount tx data(invalid currency), should return error except update domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   true,
			sendResult:     false,
		}
		setupForTx(tr)
		ctx = assemblyCtxData("OTT", int64(18), false, false, false, nil, false, false, 0, true, 1)
		for txType, txRefer := range txTypeList {
			tx, _ := assemblyTxData(false, "100", txType, false, false, true)
			ok, _ := txRefer.tx.Validate(ctx, tx)
			assert.Equal(t, txRefer.want, ok, "%s domain Validate fails ", txType)
		}
	})
	t.Run("not enough token, should return error for create domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: true,
			saleResult:     true,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		ctx = assemblyCtxData("OLT", int64(18), false, false, false, nil, false, false, 0, true, 1)
		for txType, txRefer := range txTypeList {
			tx, _ := assemblyTxData(false, "10", txType, false, false, true)
			ok, _ := txRefer.tx.Validate(ctx, tx)
			assert.Equal(t, txRefer.want, ok, "%s domain Validate fails ", txType)
		}
	})
	t.Run("valid test data, should return ok", func(t *testing.T) {
		tr := &txResult{
			createResult:   true,
			purchaseResult: true,
			saleResult:     true,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100000000000000", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(10), true, false, true, owner, false, false, 0, true, 1)
			ok, _ := txRefer.tx.Validate(ctx, tx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain Validate fails ", txType)
		}
	})
}

func TestDomainTx_ProcessCheck(t *testing.T) {
	ctx := &action.Context{}
	//t.Run("cast tx with invalid tx data, should return error for each tx type", func(t *testing.T) {
	//	tr := &txResult{
	//		createResult:   false,
	//		purchaseResult: false,
	//		saleResult:     false,
	//		updateResult:   false,
	//		sendResult:     false,
	//	}
	//	setupForTx(tr)
	//	for txType, txRefer := range txTypeList {
	//		tx, _ := assemblyTxData(false, "100", txType, true, false, true)
	//		ctx = assemblyCtxData("", int64(10), false, true, false, nil, false, false, 0, true, 1)
	//		ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
	//		assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
	//	}
	//})
	t.Run("get balance, domain, sender balance with invalid data, should return error", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, _ := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(10), true, true, false, nil, true, false, 0, true, 1)
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("non-exist domain, should return error except create domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   true,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(10), true, true, true, owner, true, false, 0, true, 1)
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("invalid currency type OLL, non-exist domain, should return error", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLL", int64(10), true, true, true, owner, true, false, 0, true, 1)
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("create a domain, should return error except send to domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     true,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(10), true, true, true, owner, true, false, 0, true, 1)
			// create a test domain before test
			d := &ons.Domain{
				Name: "test_domain",
			}
			err := ctx.Domains.Set(d)
			ctx.Domains.State.Commit()
			assert.Nil(t, err, "pre setup a test domain, should be no error")
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("create a domain with some keypair, should return error except send to domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		owner, ownerPubkey, ownerPrikey := generateKeyPair()
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx := assemblyTxDataSameKp(false, "100", txType, false, owner, ownerPubkey, ownerPrikey)
			ctx = assemblyCtxData("OLT", int64(10), true, true, true, owner, true, false, 0, true, 1)
			// create a test domain before test
			d := &ons.Domain{
				Name:         "test_domain",
				OwnerAddress: owner.Bytes(),
			}
			err := ctx.Domains.Set(d)
			ctx.Domains.State.Commit()
			assert.Nil(t, err, "pre setup a test domain, should be no error")
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("create a domain and set it onsale with high price(same keypair)", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		owner, ownerPubkey, ownerPrikey := generateKeyPair()
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx := assemblyTxDataSameKp(false, "1000", txType, false, owner, ownerPubkey, ownerPrikey)
			ctx = assemblyCtxData("OLT", int64(5), true, true, true, owner, true, true, 1000000000000000000, true, 1)
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("create a domain and set it onsale with reasonable price(same keypair)", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: true,
			saleResult:     false,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		owner, ownerPubkey, ownerPrikey := generateKeyPair()
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx := assemblyTxDataSameKp(false, "1000", txType, false, owner, ownerPubkey, ownerPrikey)
			ctx = assemblyCtxData("OLT", int64(5), true, true, true, owner, true, true, 100, true, 1)
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
	t.Run("create a domain and set it onsale with high price,reset ctx header height", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, true, true)
			ctx = assemblyCtxData("OLT", int64(5), true, true, true, owner, true, true, 10000, false, 1000)
			ok, _ := txRefer.tx.ProcessCheck(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessCheck fails ", txType)
		}
	})
}

func TestDomainTx_ProcessDeliver(t *testing.T) {
	ctx := &action.Context{}
	//t.Run("cast tx with invalid tx data, should return error for each tx type", func(t *testing.T) {
	//	tr := &txResult{
	//		createResult:   false,
	//		purchaseResult: false,
	//		saleResult:     false,
	//		updateResult:   false,
	//		sendResult:     false,
	//	}
	//	setupForTx(tr)
	//	for txType, txRefer := range txTypeList {
	//		tx, _ := assemblyTxData(false, "100", txType, true, false, true)
	//		ctx = assemblyCtxData("", int64(10), false, true, false, nil, false, false, 0, true, 1)
	//		ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
	//		assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
	//	}
	//})
	t.Run("get balance, domain, sender balance with invalid data, should return error", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, _ := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(10), true, true, false, nil, true, false, 0, true, 1)
			ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
		}
	})
	t.Run("non-exist domain, should return error except create domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   true,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(10), true, true, true, owner, true, false, 0, true, 1)
			ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
		}
	})
	t.Run("invalid currency type OLL, non-exist domain", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     false,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLL", int64(10), true, true, true, owner, true, false, 0, true, 1)
			ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
		}
	})
	t.Run("create a domain with high price and false with onsale", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: false,
			saleResult:     false,
			updateResult:   false,
			sendResult:     true,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "100", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(5), true, true, true, owner, true, true, 10000000000000, false, 1)
			ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
		}
	})
	t.Run("create a domain with high price and true with onsale", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: true,
			saleResult:     false,
			updateResult:   false,
			sendResult:     true,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "10000000000000", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(5), true, true, true, owner, true, true, 100000, true, 1)
			ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
		}
	})
	t.Run("create a domain with reasonable price and true with onsale, reset ctx header height", func(t *testing.T) {
		tr := &txResult{
			createResult:   false,
			purchaseResult: true,
			saleResult:     true,
			updateResult:   true,
			sendResult:     true,
		}
		setupForTx(tr)
		for txType, txRefer := range txTypeList {
			testDB := setup()
			tx, owner := assemblyTxData(false, "10000", txType, false, false, true)
			ctx = assemblyCtxData("OLT", int64(5), true, true, true, owner, true, true, 100, true, 1000)
			ok, _ := txRefer.tx.ProcessDeliver(ctx, tx.RawTx)
			teardown(testDB)
			assert.Equal(t, txRefer.want, ok, "%s domain ProcessDeliver fails ", txType)
		}
	})
}
