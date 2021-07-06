package eth

import (
	"github.com/Oneledger/protocol/data/accounts"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// Accounts returns the list of accounts available to this node.
func (svc *Service) Accounts() ([]common.Address, error) {
	svc.logger.Debug("eth_accounts")
	svc.mu.Lock()
	defer svc.mu.Unlock()

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

// GetBalance returns the provided account's balance up to the provided block number.
func (svc *Service) GetBalance(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		return &hexutil.Big{}, err
	}
	svc.logger.Debug("eth_getBalance", "height", height)

	ethAcc, err := svc.ctx.GetAccountKeeper().GetVersionedAccount(svc.getStateHeight(height), address.Bytes())
	if err != nil {
		svc.logger.Debug("eth_getBalance", "account_not_found", address)
		return (*hexutil.Big)(ethAcc.Balance()), nil
	}
	return (*hexutil.Big)(ethAcc.Balance()), nil
}
