package utils

import "strings"

func PadZero(s string) string {
	ss := strings.Split(s, ".")
	if len(ss) == 2 {
		ss = []string{strings.TrimLeft(ss[0], "0"), strings.TrimLeft(ss[1], "0"), strings.Repeat("0", 18-len(ss[1]))}
	} else {
		ss = []string{strings.TrimLeft(ss[0], "0"), strings.Repeat("0", 18)}
	}
	s = strings.Join(ss, "")
	return s
}
