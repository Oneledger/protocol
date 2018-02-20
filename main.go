package main

import "./core"
import "fmt"
import "strconv"

func main() {
  bc := core.NewBlockchain()
  bc.AddBlock("Send 1 BTC to YangLi")
  bc.AddBlock("Send 2 BTC to YangLi")

  for _, block := range bc.Blocks {
    fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
    fmt.Printf("Data: %s\n", block.Data)
    fmt.Printf("Hash: %x\n",block.Hash)
    pow := core.NewProofOfWork(block)
    fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
    fmt.Println()
  }
}
