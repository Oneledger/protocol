package ons

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainUpdate{}

type DomainUpdate struct {
	Owner       action.Address `json:"owner"`
	Beneficiary action.Address `json:"beneficiary"`
	Name        ons.Name       `json:"name"`
	Active      bool           `json:"active"`
	Uri         string         `json:"uri"`
}

func (du DomainUpdate) Marshal() ([]byte, error) {
	return json.Marshal(du)
}

func (du *DomainUpdate) Unmarshal(data []byte) error {
	return json.Unmarshal(data, du)
}

func (du DomainUpdate) OnsName() string {
	return du.Name.String()
}

func (du DomainUpdate) Signers() []action.Address {
	return []action.Address{du.Owner}
}

func (du DomainUpdate) Type() action.Type {
	return action.DOMAIN_UPDATE
}

func (du DomainUpdate) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(du.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.from"),
		Value: du.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = domainUpdateTx{}

type domainUpdateTx struct {
}

func (domainUpdateTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	update := &DomainUpdate{}
	err := update.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(tx.RawBytes(), update.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if update.Owner == nil || len(update.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !update.Name.IsValid() {
		return false, ErrInvalidDomain
	}

	if update.Active == false && update.Beneficiary == nil {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainUpdateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runUpdate(ctx, tx)
}

func (domainUpdateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runUpdate(ctx, tx)
}

func (domainUpdateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runUpdate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	update := &DomainUpdate{}
	err := update.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	if !ctx.Domains.Exists(update.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain doesn't exist: %s", update.Name)}
	}

	d, err := ctx.Domains.Get(update.Name)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get domain: %s", update.Name)}
	}

	if !d.IsChangeable(ctx.Header.Height) {
		return false, action.Response{Log: fmt.Sprintf("domain is not changable: %s, last change: %d ,current height :%d", update.Name, d.LastUpdateHeight, ctx.Header.Height)}
	}

	if !bytes.Equal(d.Owner, update.Owner) {
		return false, action.Response{Log: fmt.Sprintf("domain is not owned by: %s", hex.EncodeToString(update.Owner))}
	}

	d.SetAccountAddress(update.Beneficiary)
	if update.Active {
		d.Activate()
	} else {
		d.Deactivate()

		// deactivate all subdomains
		if !d.Name.IsSub() {

			ctx.Domains.IterateSubDomain(d.Name, func(name ons.Name, domain *ons.Domain) bool {
				domain.Deactivate()
				err := ctx.Domains.Set(domain)
				if err != nil {
					ctx.Logger.Error("failed to update sub domain activate status ", domain.Name, err)
					return false
				}

				d.SetLastUpdatedHeight(ctx.Header.Height)
				return false
			})
		}
	}

	d.SetLastUpdatedHeight(ctx.Header.Height)
	opt, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, gov.ErrGetONSOptions, update.Tags(), err)
	}
	if len(update.Uri) > 0 {
		ok := opt.IsValidURI(update.Uri)
		if !ok {
			return false, action.Response{
				Log: "invalid uri provided",
			}
		}
		d.URI = update.Uri
	} else if strings.Compare(update.Uri, "") == 0 {
		// reset uri
		d.URI = ""
	}

	err = ctx.Domains.Set(d)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	return true, action.Response{Events: action.GetEvent(update.Tags(), "update_domain")}
}
