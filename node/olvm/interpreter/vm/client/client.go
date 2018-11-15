package client

import (
  "net/rpc"
  service ".."
)

type OLVMClient struct {
  Protocol string
  ServicePath string
}

var DefaultClient = MakeClient("tcp", "localhost:1980")

func (c OLVMClient) Run(from string, address string, callString string, value int) (service.Reply, error){
  args := service.Args{from, address, callString, value}
  var reply service.Reply
  client, err := rpc.DialHTTP(c.Protocol, c.ServicePath)
  if err != nil {
    return reply, err
  }
  err = client.Call("Container.Exec", &args, &reply)
  if err != nil {
    return reply, err
  }

  return reply, nil
}

func MakeClient(protocol, path string) OLVMClient{
  return OLVMClient{protocol, path}
}

func Run(from string, address string, callString string, value int) (service.Reply, error) {
  return DefaultClient.Run(from, address, callString, value)
}
