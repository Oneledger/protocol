package farm_action

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"
)

func GetProductStore(ctx *action.Context) (*farm_data.ProductStore, error) {
	store, err := ctx.ExtStores.Get("extProductStore")
	if err != nil {
		return nil, err
	}
	productStore, ok := store.(*farm_data.ProductStore)
	if ok == false {
		return nil, err
	}

	return productStore, nil
}
