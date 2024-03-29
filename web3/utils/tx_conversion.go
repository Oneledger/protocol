package utils

import (
	"math/big"
	"strconv"

	"github.com/Oneledger/protocol/action"
	olvm "github.com/Oneledger/protocol/action/olvm"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	jsonSerializer = serialize.GetSerializer(serialize.NETWORK)
	big8           = big.NewInt(8)
)

func EthToOLSignedTx(tx *ethtypes.Transaction) (*action.SignedTx, error) {
	chainId := tx.ChainId()

	signer := ethtypes.NewEIP155Signer(chainId)
	chainIdMul := new(big.Int).Mul(chainId, big.NewInt(2))
	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, chainIdMul)
	V.Sub(V, big8)

	sighash := signer.Hash(tx)
	pub, err := utils.RecoverPlain(sighash, R, S, V, true)
	if err != nil {
		return nil, err
	}
	pubKey, err := keys.GetPublicKeyFromBytes(crypto.CompressPubkey(pub), keys.ETHSECP)
	if err != nil {
		return nil, err
	}

	handler, err := pubKey.GetHandler()
	if err != nil {
		return nil, err
	}

	sigs := []action.Signature{{Signer: pubKey, Signed: utils.ToUncompressedSig(R, S, V)}}

	// fee setting
	gas := int64(tx.Gas())
	gasPrice := tx.GasPrice()
	price := action.NewAmount(action.DEFAULT_CURRENCY, balance.Amount(*gasPrice))
	fee := action.Fee{Price: *price, Gas: gas}

	// data conversion
	value := tx.Value()
	var to *action.Address
	if tx.To() != nil {
		to = new(action.Address)
		*to = tx.To().Bytes()
	}
	sendAmount := action.NewAmount(action.DEFAULT_CURRENCY, balance.Amount(*value))
	msg := olvm.Transaction{
		From:    handler.Address(),
		To:      to,
		Amount:  *sendAmount,
		Data:    tx.Data(),
		Nonce:   tx.Nonce(),
		ChainID: tx.ChainId(),
	}
	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	rawTx := action.RawTx{
		Type: action.OLVM,
		Data: data,
		Fee:  fee,
		Memo: strconv.FormatUint(tx.Nonce(), 10),
	}
	signedTx := action.SignedTx{
		RawTx:      rawTx,
		Signatures: sigs,
	}
	return &signedTx, nil
}
