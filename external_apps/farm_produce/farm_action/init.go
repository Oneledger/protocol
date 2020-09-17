package farm_action

import "github.com/Oneledger/protocol/action"

const (
	FARM_INSERT_PRODUCT action.Type = 990101
)

func init() {
	action.RegisterTxType(FARM_INSERT_PRODUCT, "FARM_INSERT_PRODUCT")
}