package main

import (
  "log"
  "net/rpc"
  "./service"
)

func main () {
  client, err := rpc.DialHTTP("tcp", "localhost:1980")
  if err!= nil {
    log.Fatal("dialing:", err)
  }

  args := service.Args{}
  args.Address = "samples://deadloop"

  var reply service.Reply

  err = client.Call("OLVMService.Exec", &args, &reply)

  if err != nil {
    log.Fatal("service error:", err)
  }
  log.Printf("result: %v", reply)
}
