package service

import (
  "log"
  "net"
  "net/rpc"
  "net/http"
  "os"
)

func (ol *OLVMService) Echo(args *Args, reply *int) error {
  *reply = 200
  return nil
}

func (ol *OLVMService) Close(args *Args, reply *int) error {
  defer func() {
    os.Exit(0)
  }()
  *reply = 200
  return nil
}

func Run() {
  log.Print("Service is running...")
  service := new(OLVMService)
  rpc.Register(service)
  rpc.HandleHTTP()
  l, e := net.Listen("tcp", ":1980")
  if e!= nil {
    log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)

}
