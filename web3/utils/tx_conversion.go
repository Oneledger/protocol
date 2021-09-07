package utils

import (
	"math/big"

	"github.com/Oneledger/protocol/action"
	olvm "github.com/Oneledger/protocol/action/olvm"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	jsonSerializer = serialize.GetSerializer(serialize.NETWORK)
	big8           = big.NewInt(8)
)

func EthToOLSignedTx(tx *ethtypes.Transaction) (tmtypes.Tx, error) {
	chainId := tx.ChainId()

	signer := ethtypes.NewEIP155Signer(chainId)
	chainIdMul := new(big.Int).Mul(chainId, big.NewInt(2))
	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, chainIdMul)
	V.Sub(V, big8)

	sighash := signer.Hash(tx)
	pub, err := utils.RecoverPlain(sighash, R, S, V, true)
	if err != nil {
		return []byte{}, err
	}
	pubKey, err := keys.GetPublicKeyFromBytes(crypto.CompressPubkey(pub), keys.ETHSECP)
	if err != nil {
		return []byte{}, err
	}

	handler, err := pubKey.GetHandler()
	if err != nil {
		return []byte{}, err
	}

	sigs := []action.Signature{{Signer: pubKey, Signed: utils.ToUncompressedSig(R, S, V)}}

	// memo creation
	uuidNew, err := uuid.NewUUID()
	if err != nil {
		return []byte{}, err
	}
	memo := uuidNew.String()

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
	msg := olvm.Nexus{
		From:    handler.Address(),
		To:      to,
		Amount:  *sendAmount,
		Data:    tx.Data(),
		Nonce:   tx.Nonce(),
		ChainID: tx.ChainId(),
	}
	data, err := msg.Marshal()
	if err != nil {
		return []byte{}, err
	}

	rawTx := action.RawTx{
		Type: action.NEXUS,
		Data: data,
		Fee:  fee,
		Memo: memo,
	}
	signedTx := action.SignedTx{
		RawTx:      rawTx,
		Signatures: sigs,
	}

	packet, err := jsonSerializer.Serialize(signedTx)
	if err != nil {
		return []byte{}, err
	}
	return packet, nil
}
