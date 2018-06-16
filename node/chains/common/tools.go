/*
	Copyright 2017 - 2018 OneLedger
*/
package common

import (
	"io"
	"net/http"
)

type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type Logger interface {
	Println(v ...interface{})
}
