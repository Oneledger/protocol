package vpart_action

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_data"
)

func GetVPartStore(ctx *action.Context) (*vpart_data.VPartStore, error) {
	store, err := ctx.ExtStores.Get("extVPartStore")
	if err != nil {
		return nil, err
	}
	produceStore, ok := store.(*vpart_data.VPartStore)
	if ok == false {
		return nil, err
	}

	return produceStore, nil
}
