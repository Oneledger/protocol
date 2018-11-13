package main

import (
  "log"
  "./service"
)

func main() {
  log.Print("Up running vm")
  service.Run()
}
