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

  var reply int

  err = client.Call("OLVMService.Close", &args, &reply)

  if err != nil {
    log.Fatal("service error:", err)
  }
  log.Printf("result: %d", reply)
}
