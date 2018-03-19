package types

import "net/url"

type Reference struct {
	Type string `json:"type"`
	Url url.URL `json:"type"`
}