package farm_action

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"
)

func GetProduceStore(ctx *action.Context) (*farm_data.ProduceStore, error) {
	store, err := ctx.ExtStores.Get("extProduceStore")
	if err != nil {
		return nil, err
	}
	productStore, ok := store.(*farm_data.ProduceStore)
	if ok == false {
		return nil, err
	}

	return productStore, nil
}
