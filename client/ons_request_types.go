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
	Owner   keys.Address  `json:"owner"`
	Account keys.Address  `json:"account"`
	Name    string        `json:"name"`
	Price   action.Amount `json:"price"`
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type ONSUpdateRequest struct {
	Owner   keys.Address  `json:"owner"`
	Account keys.Address  `json:"account"`
	Name    string        `json:"name"`
	Active  bool          `json:"active"`
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type ONSSaleRequest struct {
	Name         string        `json:"name"`
	OwnerAddress keys.Address  `json:"owner"`
	Price        action.Amount `json:"price"`
	CancelSale   bool          `json:"cancel_sale"`
	Fee          action.Amount `json:"fee"`
	Gas          int64         `json:"gas"`
}

type ONSPurchaseRequest struct {
	Name     string        `json:"name"`
	Buyer    keys.Address  `json:"buyer"`
	Account  keys.Address  `json:"account"`
	Offering action.Amount `json:"offering"`
	Fee      action.Amount `json:"fee"`
	Gas      int64         `json:"gas"`
}

type ONSSendRequest struct {
	From   keys.Address  `json:"from"`
	Name   string        `json:"name"`
	Amount action.Amount `json:"amount"`
	Fee    action.Amount `json:"fee"`
	Gas    int64         `json:"gas"`
}

type ONSGetDomainsRequest struct {
	Name   string       `json:"name"`
	Owner  keys.Address `json:"owner"`
	OnSale bool         `json:"onSale"`
}

type ONSGetDomainsReplay struct {
	Domains []ons.Domain `json:"domains"`
}
