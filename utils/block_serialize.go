package utils

func Serialize(b *Block) []byte {
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
