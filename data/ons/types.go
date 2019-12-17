package ons

import (
	"regexp"
)

const (
	reg    = `^([a-zA-Z0-9]+\.)*[a-zA-Z0-9]+\.[a-zA-Z]{2,11}?$`
	sub    = `^([a-zA-Z0-9]+\.)+[a-zA-Z0-9]+\.[a-zA-Z]{2,11}?$`
	parent = `[a-zA-Z0-9]+\.[a-zA-Z]{2,11}?$`
)

var (
	pattern    *regexp.Regexp
	subpattern *regexp.Regexp
)

func init() {
	pattern = regexp.MustCompile(reg)
	subpattern = regexp.MustCompile(sub)
}

type Name string

func GetNameFromString(s string) Name {
	return Name(s)
}

func (n Name) String() string {
	return string(n)
}

func (n Name) IsValid() bool {
	if len(n) > 256 {
		return false
	}
	return pattern.Match([]byte(n.String()))
}

func (n Name) IsSub() bool {
	return subpattern.Match([]byte(n.String()))
}

func (n Name) toKey() []byte {
	return []byte(reverse(n.String()))
}

func getIndex(s string) (int, int) {
	a, b := 0, 0
	for i, c := range []rune(s) {
		if c == rune('.') {
			a = b
			b = i
		}
	}
	return a, b
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}
