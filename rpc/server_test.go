package rpc

import (
	"io/ioutil"
	"net/rpc"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

// Simple round-trip test
func TestServer(t *testing.T) {
	srv := NewServer(ioutil.Discard)
	defer srv.Close()
	assert.NotNil(t, srv.rpc)
	assert.NotNil(t, srv.http)
	assert.NotNil(t, srv.logger)
	assert.NotNil(t, srv.mux)

	service := Arith(0)
	u, _ := url.Parse("tcp://127.0.0.1:9006")
	err := srv.Prepare(u, &service)
	assert.Nil(t, err, "preparing the server shouldn't return an error")

	err = srv.Start()
	assert.Nil(t, err, "start should be ok")

	time.Sleep(1 * time.Second)

	client, err := rpc.DialHTTPPath(u.Scheme, u.Host, RPCPath)
	if err != nil {
		assert.Nil(t, err, "failed to connect to server: %s", err.Error())
	}

	var reply int
	args := Args{3, 5}
	err = client.Call("Arith.Multiply", args, &reply)
	assert.Nil(t, err, "got an error calling the server")

	assert.Equal(t, 15, reply, "service should calculate the right result")
}
