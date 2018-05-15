package common

import (
	"net/http"
	"io"
)


type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type Logger interface {
	Println(v ...interface{})
}