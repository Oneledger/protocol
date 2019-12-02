package event

import (
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/utils/transition"
)

func init() {

	err := EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.SIGNING,
		Fn:   Signing,
		From: transition.Status(ethereum.New),
		To:   transition.Status(ethereum.BusyBroadcasting),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.FINALIZESIGNING,
		Fn:   FinalizeSigning,
		From: transition.Status(ethereum.BusyBroadcasting),
		To:   transition.Status(ethereum.BusyFinalizing),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.VERIFYREDEEM,
		Fn:   VerifyRedeem,
		From: transition.Status(ethereum.BusyFinalizing),
		To:   transition.Status(ethereum.Finalized),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.BURN,
		Fn:   Burn,
		From: transition.Status(ethereum.Finalized),
		To:   transition.Status(ethereum.Released),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.CLEANUP,
		Fn:   redeemcleanup,
		From: transition.Status(ethereum.Released),
		To:   transition.Status(0),
	})
	if err != nil {
		panic(err)
	}
}

func Signing(ctx interface{}) error {
	return nil
}

func FinalizeSigning(ctx interface{}) error {
	return nil
}

func VerifyRedeem(ctx interface{}) error {
	return nil
}

func Burn(ctx interface{}) error {
	return nil
}

func redeemcleanup(ctx interface{}) error {
	return nil
}
