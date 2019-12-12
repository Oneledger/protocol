package ons

import (
	"strings"
)

type Name string

func GetNameFromString(s string) Name {

	return Name(s)
}

func (n Name) String() string {
	return string(n)
}

func (n Name) IsValid() bool {
	arr := strings.Split(n.String(), ".")
	if len(arr) < 2 {
		return false
	}
	return true
}
