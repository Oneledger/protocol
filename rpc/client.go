package rpc

import (
	"net/rpc"
	"net/url"
)

type Client struct {
	*rpc.Client
}

func NewClient(addr string) (*Client, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	// For gobs and http
	client, err := rpc.Dial(u.Scheme, u.Host)
	if err != nil {
		return nil, err
	}

	// TODO: for jsonrpc 2.0
	//conn, err := net.Dial(u.Scheme, u.Host)
	//if err != nil {
	//	return nil, err
	//}

	//codec := jsonrpc2.NewClientCodec(conn)
	//client := rpc.NewClientWithCodec(codec)
	return &Client{client}, nil
}
