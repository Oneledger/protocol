package core

type Blockchain struct {
  Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
  prevBlock := bc.Blocks[len(bc.Blocks)-1]
  newBlock := NewBlock(data, prevBlock.Hash)
  bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block {
  return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
  return &Blockchain{[]*Block{NewGenesisBlock()}}
}
