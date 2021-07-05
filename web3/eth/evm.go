package eth

import (
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

func (svc *Service) GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	contractStore := svc.stateDB.GetContractStore()

	height, err := StateAndHeaderByNumberOrHash(svc.ext.RPCClient(), blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	switch height {
	case -1:
	case -2:
		height = contractStore.State.Version()
		break
	}
	prefixKey := utils.GetStorageByAddressKey(address, common.HexToHash(key).Bytes())
	prefixStore := evm.AddressStoragePrefix(address)

	value := common.Hash{}
	storeKey := contractStore.GetStoreKey(prefixStore, prefixKey.Bytes())
	rawValue := contractStore.State.GetVersioned(height, storeKey)
	if len(rawValue) > 0 {
		value.SetBytes(rawValue)
	}
	return value[:], nil
}

func (svc *Service) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	accountKeeper := svc.stateDB.GetAccountKeeper()
	contractStore := svc.stateDB.GetContractStore()

	height, err := StateAndHeaderByNumberOrHash(svc.ext.RPCClient(), blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	switch height {
	case -1:
	case -2:
		height = contractStore.State.Version()
		break
	}

	ethAcc, err := accountKeeper.GetVersionedAccount(height, address.Bytes())
	if err != nil {
		return hexutil.Bytes{}, nil
	}
	storeKey := contractStore.GetStoreKey(evm.KeyPrefixCode, ethAcc.CodeHash)
	code := contractStore.State.GetVersioned(height, storeKey)
	return code[:], nil
}
