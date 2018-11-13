package service

import (
  "log"
  "net"
  "net/rpc"
  "net/http"
  "os"
  "../monitor"
  "../runner"
  "errors"
  "fmt"
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

func (ol *OLVMService) Exec(args *Args, reply *Reply) error {
  mo := monitor.CreateMonitor(10, monitor.DEFAULT_MODE, "./ovm.pid")
  status_ch := make(chan monitor.Status)
  runner_ch := make(chan bool)
  go mo.CheckStatus(status_ch)
  go func () {
    runner := runner.CreateRunner()
    from, address, callString, value := parseArgs(args)
    out, ret := runner.Call(from, address, callString, value)
    reply.Out = out
    reply.Ret = ret
    runner_ch <- true
  }()
  for {
    select {
    case <- runner_ch:
      return nil
    case status := <- status_ch:
      panic(status)
      return errors.New(fmt.Sprintf("%s : %d", status.Details, status.Code))
    }
  }

  return nil
}

func parseArgs(args *Args) (string, string, string, int){
  from := args.From
  address := args.Address
  callString := args.CallString
  value := args.Value
  if from == "" {
    from = "0x0"
  }
  if address == "" {
    address = "samples://helloworld"
  }
  if callString == "" {
    callString = "default__('hello,world from Oneledger')"
  }
  return from, address, callString, value
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
