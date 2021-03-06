package network_delegation

// below is removed since withdraw logic is moved to block beginner, OLP-1267

//var _ action.Msg = &WithdrawNetworkDelegation{}
//
//type WithdrawNetworkDelegation struct {
//	DelegationAddress keys.Address
//	Amount            action.Amount
//}
//
//func (w WithdrawNetworkDelegation) Signers() []action.Address {
//	return []action.Address{w.DelegationAddress}
//}
//
//func (w WithdrawNetworkDelegation) Type() action.Type {
//	return action.WITHDRAW_NETWORK_DELEGATION
//}
//
//func (w WithdrawNetworkDelegation) Tags() kv.Pairs {
//	tags := make([]kv.Pair, 0)
//
//	tag1 := kv.Pair{
//		Key:   []byte("tx.type"),
//		Value: []byte(w.Type().String()),
//	}
//	tag2 := kv.Pair{
//		Key:   []byte("tx.delegationAddress"),
//		Value: w.DelegationAddress.Bytes(),
//	}
//
//	tags = append(tags, tag1, tag2)
//	return tags
//}
//
//func (w WithdrawNetworkDelegation) Marshal() ([]byte, error) {
//	return json.Marshal(w)
//}
//
//func (w *WithdrawNetworkDelegation) Unmarshal(bytes []byte) error {
//	return json.Unmarshal(bytes, w)
//}
//
//var _ action.Tx = withdrawNetworkDelegationTx{}
//
//type withdrawNetworkDelegationTx struct{}
//
//func (n withdrawNetworkDelegationTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
//	withdraw := &WithdrawNetworkDelegation{}
//	err := withdraw.Unmarshal(tx.Data)
//	if err != nil {
//		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
//	}
//	err = action.ValidateBasic(tx.RawBytes(), withdraw.Signers(), tx.Signatures)
//	if err != nil {
//		return false, err
//	}
//
//	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
//	if err != nil {
//		return false, err
//	}
//
//	if err := withdraw.DelegationAddress.Err(); err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
//
//func (n withdrawNetworkDelegationTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	return runWithdrawNetworkDelegation(ctx, tx)
//}
//
//func (n withdrawNetworkDelegationTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	return runWithdrawNetworkDelegation(ctx, tx)
//}
//
//func (n withdrawNetworkDelegationTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
//	ctx.Logger.Detail("Processing Delegate Transaction for ProcessFee", signedTx)
//	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
//}
//
//func runWithdrawNetworkDelegation(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	withdraw := WithdrawNetworkDelegation{}
//	err := withdraw.Unmarshal(tx.Data)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, withdraw.Tags(), err)
//	}
//
//	// Check if withdraw amount is valid
//	coin := withdraw.Amount.ToCoinWithBase(ctx.Currencies)
//	if !coin.IsValid() {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidAmount, withdraw.Tags(), errors.New("Coin is not valid"))
//	}
//	if coin.Currency.Name != "OLT" {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, withdraw.Tags(), errors.New("currency is not OLT"))
//	}
//	_, ok := ctx.Currencies.GetCurrencyByName("OLT")
//	if !ok {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, withdraw.Tags(), errors.New("currency OLT does not exist in system"))
//	}
//
//	//Check matured balance
//	maturedCoin, err := ctx.NetwkDelegators.Deleg.WithPrefix(network_delegation.MatureType).Get(withdraw.DelegationAddress)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidAddress, withdraw.Tags(), err)
//	}
//	if !coin.LessThanEqualCoin(*maturedCoin) {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrNotEnoughFund, withdraw.Tags(), err)
//	}
//	// Reduce the balance withdraw amount
//	newCoin, err := maturedCoin.Minus(coin)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, withdraw.Tags(), err)
//	}
//	err = ctx.NetwkDelegators.Deleg.WithPrefix(network_delegation.MatureType).Set(withdraw.DelegationAddress, &newCoin)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, withdraw.Tags(), err)
//	}
//
//	//Add it to users address
//	err = ctx.Balances.AddToAddress(withdraw.DelegationAddress, coin)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, withdraw.Tags(), err)
//	}
//
//	return helpers.LogAndReturnTrue(ctx.Logger, withdraw.Tags(), "delegation withdraw success")
//}
