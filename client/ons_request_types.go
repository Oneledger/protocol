package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
)

/*
	ONS Request Types
*/
type ONSCreateRequest struct {
	Owner       keys.Address  `json:"owner"`
	Account     keys.Address  `json:"account"`
	Name        string        `json:"name"`
	Uri         string        `json:"uri"`
	BuyingPrice action.Amount `json:"buyingPrice"`
	GasPrice    action.Amount `json:"gasPrice"`
	Gas         int64         `json:"gas"`
}

type ONSCreateSubRequest struct {
	Owner       keys.Address  `json:"owner"`
	Account     keys.Address  `json:"account"`
	Name        string        `json:"name"`
	Uri         string        `json:"uri"`
	BuyingPrice action.Amount `json:"buyingPrice"`
	GasPrice    action.Amount `json:"gasPrice"`
	Gas         int64         `json:"gas"`
}

type ONSUpdateRequest struct {
	Owner    keys.Address  `json:"owner"`
	Account  keys.Address  `json:"account"`
	Name     string        `json:"name"`
	Active   bool          `json:"active"`
	Uri      string        `json:"uri"`
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}

type ONSRenewRequest struct {
	Owner       keys.Address  `json:"owner"`
	Account     keys.Address  `json:"account"`
	Name        string        `json:"name"`
	BuyingPrice action.Amount `json:"buyingPrice"`
	GasPrice    action.Amount `json:"gasPrice"`
	Gas         int64         `json:"gas"`
}

type ONSSaleRequest struct {
	Name         string        `json:"name"`
	OwnerAddress keys.Address  `json:"owner"`
	Price        action.Amount `json:"price"`
	CancelSale   bool          `json:"cancelSale"`
	GasPrice     action.Amount `json:"gasPrice"`
	Gas          int64         `json:"gas"`
}

type ONSPurchaseRequest struct {
	Name     string        `json:"name"`
	Buyer    keys.Address  `json:"buyer"`
	Account  keys.Address  `json:"account"`
	Offering action.Amount `json:"offering"`
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}

type ONSSendRequest struct {
	From     keys.Address  `json:"from"`
	Name     string        `json:"name"`
	Amount   action.Amount `json:"amount"`
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}

type ONSDeleteSubRequest struct {
	Name     string        `json:"name"`
	Owner    keys.Address  `json:"owner"`
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}

type ONSGetDomainsRequest struct {
	Name        string       `json:"name"`
	Owner       keys.Address `json:"owner"`
	OnSale      bool         `json:"onSale"`
	Beneficiary keys.Address `json:"beneficiary"`
}

type ONSGetDomainsReply struct {
	Domains []ons.Domain `json:"domains"`
}

type ONSGetDomainsOnSaleReply struct {
	Domains []ons.Domain `json:"domains"`
}

type ONSGetOptionsReply struct {
	ons.Options `json:"options"`
}
