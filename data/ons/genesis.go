package ons

import (
	"net/url"
	"strings"

	"github.com/Oneledger/protocol/data/balance"
)

type Options struct {
	Currency          string         `json:"currency"`
	PerBlockFees      balance.Amount `json:"perBlockFees"`
	BaseDomainPrice   balance.Amount `json:"baseDomainPrice"`
	FirstLevelDomains []string       `json:"firstLevelDomains"`

	firstLevel map[string]bool
	protocols  map[string]bool
}

func (opt *Options) IsNameAllowed(name Name) bool {
	if opt.firstLevel == nil {

		opt.firstLevel = make(map[string]bool)
		for i := 0; i < len(opt.FirstLevelDomains); i++ {
			opt.firstLevel[opt.FirstLevelDomains[i]] = true
		}
	}

	nameAr := strings.Split(name.String(), ".")
	_, ok := opt.firstLevel[nameAr[len(nameAr)-1]]
	return ok
}

func (opt *Options) IsValidURI(uri string) bool {
	if opt.protocols == nil {
		opt.protocols = map[string]bool{
			"http":  true,
			"https": true,
			"ipfs":  true,
			"ftp":   true,
		}
	}

	u, err := url.Parse(uri)
	if err != nil {
		return false
	}

	_, ok := opt.protocols[u.Scheme]
	return ok
}
