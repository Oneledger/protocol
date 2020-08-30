/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package main

import "strings"

func padZero(s string) string {
	ss := strings.Split(s, ".")
	if len(ss) == 2 {
		ss = []string{strings.TrimLeft(ss[0], "0"), strings.TrimLeft(ss[1], "0"), strings.Repeat("0", 18-len(ss[1]))}
	} else {
		ss = []string{strings.TrimLeft(ss[0], "0"), strings.Repeat("0", 18)}
	}
	s = strings.Join(ss, "")
	return s
}
