package evm

// func NewEVM(ctx *Context, sv *config.Server, txCtx *vm.TxContext) *vm.EVM {
// 	blockCtx := vm.BlockContext{
// 		CanTransfer: core.CanTransfer,
// 		Transfer:    core.Transfer,
// 		GetHash:     nil,
// 		Coinbase:    common.Address{}, // there's no beneficiary since we're not mining
// 		GasLimit:    672197400,        // TODO: Set gas limit from the settings
// 		BlockNumber: big.NewInt(ctx.Header.GetHeight()),
// 		Time:        big.NewInt(ctx.Header.Time.Unix()),
// 		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
// 	}

// 	// TODO: Add EIPS
// 	eips := make([]int, 0)
// 	// for i, eip := range extraEIPs {
// 	// 	eips[i] = int(eip)
// 	// }

// 	vmConfig := vm.Config{
// 		ExtraEips: eips,
// 	}

// 	ethConfig, _ := sv.EthereumConfig()
// 	return vm.NewEVM(blockCtx, *txCtx, csdb, ethConfig, vmConfig)
// }
