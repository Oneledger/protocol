package vm

import (
  "net/rpc"
  "log"
  "strings"
  "time"
)


var DefaultClient = MakeClient("tcp", "localhost:1980")

func (c OLVMClient) Run(from string, address string, callString string, value int) (Reply, error){
  args := Args{from, address, callString, value}
  var reply Reply
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

func Run(from string, address string, callString string, value int) (Reply, error) {
  return DefaultClient.Run(from, address, callString, value)
}

func AutoRun(from string, address string, callString string, sourceCode string, value int) (reply Reply, err error) {
  reply, err = DefaultClient.Run(from, address, callString, value)
  if err!= nil && strings.HasSuffix(err.Error(),"connection refused") {
    //try to launch the service
    log.Println("trying to relaunch the service")
    go RunService()
    for err!= nil && strings.HasSuffix(err.Error(),"connection refused") {
      time.Sleep(time.Second)
      reply, err = DefaultClient.Run(from, address, callString, value)
    }
    return
  } else {
    return
  }
}
