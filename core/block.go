package core
import (
    "bytes"
    "crypto/sha256"
    "strconv"
    "time"
    "encoding/gob"
    "log"
)

type Block struct {
  Timestamp     int64
  Data          []byte
  PrevBlockHash []byte
  Hash          []byte
  Nonce         int64
}

func (b *Block) SetHash() {
  timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
  headers := bytes.Join([][]byte{b.PrevBlockHash,b.Data, timestamp},[]byte{})
  hash := sha256.Sum256(headers)
  b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
  block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{},int64(0)}
  pow := NewProofOfWork(block)
  nonce, hash := pow.Mine()
  block.Hash = hash
  block.Nonce = nonce
  return block
}

func (b *Block) Serialize() []byte {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)
  err := encoder.Encode(b)
  if err != nil {
    log.Panic(err)
  }
  return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
  var block Block
  decoder := gob.NewDecoder(bytes.NewReader(d))
  err := decoder.Decode(&block)
  if err != nil {
    log.Panic(err)
  }
  return &block
}
