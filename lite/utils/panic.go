package utils

import "log"

func RequireNil(anything interface{}) {
	if anything != nil {
		log.Panic(anything)
	}
}
