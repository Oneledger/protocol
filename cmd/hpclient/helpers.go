package main

import (
	"time"
)

func GetTime() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}
