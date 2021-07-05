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
		height = contractStore.State.Version()
		break
	case -2:
		// NOTE: Need a pending state?
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
