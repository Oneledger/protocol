/*

 */

package ons

import "github.com/pkg/errors"

var (
	ErrDomainNameNotValid = errors.New("Domain name is invalid")
	ErrDomainNotFound     = errors.New("Domain doesn't exist")
)
