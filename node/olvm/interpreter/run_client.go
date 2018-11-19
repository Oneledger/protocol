package main

import (
  "log"
  "./vm"
)

func main () {
  reply, err := vm.AutoRun("0x0","samples://helloworld","","", 0)
  if err != nil {
    log.Fatal(err)
  }
  log.Println(reply)
}
