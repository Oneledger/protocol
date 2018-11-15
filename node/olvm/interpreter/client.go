package main

import (
  "log"
  "./vm/client"
)

func main () {
  reply, err := client.Run("0x0","samples://helloworld","", 0)
  if err != nil {
    log.Fatal(err)
  }
  log.Println(reply)
}
