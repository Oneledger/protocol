package eth

import (
	"bytes"
	"math/big"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// Accounts returns the list of accounts available to this node.
func (svc *Service) Accounts() ([]common.Address, error) {
	svc.logger.Debug("eth_accounts")

	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	// NOTE: Maybe we do not need to implement this handler, but current keystore has ed25519 keys which is not
	// supported in this api, after a new structure need to think what to do with this
	wallet, err := accounts.NewWalletKeyStore("keystore/")
	if err != nil {
		return addresses, err
	}

	oltAddresses, err := wallet.ListAddresses()
	if err != nil {
		return addresses, err
	}

	for _, oltAddress := range oltAddresses {
		addresses = append(addresses, common.BytesToAddress(oltAddress))
	}

	return addresses, nil
}

func (svc *Service) getFeePoolBalance() (*big.Int, error) {
	feePool := svc.ctx.GetFeePool()
	coin, err := feePool.Get(keys.Address(fees.POOL_KEY))
	if err != nil {
		return big.NewInt(0), err
	}

	if coin.Amount == nil {
		return big.NewInt(0), nil
	}

	return coin.Amount.BigInt(), nil
}

// GetBalance returns the provided account's balance up to the provided block number.
func (svc *Service) GetBalance(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetBlockStore(), blockNrOrHash)
	if err != nil {
		svc.logger.Debug("eth_getBalance", "block err", err)
		return (*hexutil.Big)(big.NewInt(0)), nil
	}

	var (
		blockNum       int64
		pendingBalance = big.NewInt(0)
	)

	switch height {
	case rpctypes.PendingBlockNumber:
		blockNum = svc.getState().Version()
		svc.logger.Debug("eth_getBalance", "height", blockNum, "pending")
		chainID, err := svc.ChainId()
		if err != nil {
			return nil, err
		}
		rpcutils.GetPendingTxsWithCallback(svc.GetMempool(), (*big.Int)(&chainID), func(tx *rpctypes.Transaction) bool {
			if bytes.Equal(tx.From.Bytes(), address.Bytes()) {
				pendingBalance = new(big.Int).Sub(pendingBalance, tx.Value.ToInt())
			} else if tx.To != nil && bytes.Equal(tx.To.Bytes(), address.Bytes()) {
				pendingBalance = new(big.Int).Add(pendingBalance, tx.Value.ToInt())
			}
			return false
		})
	case rpctypes.LatestBlockNumber:
		blockNum = svc.getState().Version()
		svc.logger.Debug("eth_getBalance", "height", blockNum, "latest")
	case rpctypes.EarliestBlockNumber:
		blockNum = rpctypes.InitialBlockNumber
		svc.logger.Debug("eth_getBalance", "height", blockNum, "earliest")
	default:
		blockNum = height
		svc.logger.Debug("eth_getBalance", "height", blockNum)
	}

	var balance *big.Int
	acc, err := svc.ctx.GetAccountKeeper().GetVersionedAccount(address.Bytes(), blockNum)
	if err != nil {
		svc.logger.Debug("eth_getBalance", "account_not_found", address)
		balance = big.NewInt(0)
	} else {
		balance = acc.Balance()
	}

	// zero address pool
	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	if bytes.Equal(zeroAddress[:], address[:]) {
		zeroBalance, err := svc.getFeePoolBalance()
		svc.logger.Debug("eth_getBalance", "getting fee address balance")
		if err == nil {
			balance = new(big.Int).Add(balance, zeroBalance)
		}
	}

	// involve pending balance
	total := new(big.Int).Add(balance, pendingBalance)
	return (*hexutil.Big)(total), nil
}
