package core

import (
  "bytes"
  "encoding/binary"
  "log"
  "math"
  "crypto/sha256"
  "math/big"
  "fmt"
)

const targetBits = 24 //hard code for now
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
    buff := new(bytes.Buffer)
    err := binary.Write(buff, binary.BigEndian, num)
    if err != nil {
        log.Panic(err)
    }
    return buff.Bytes()
}

func NewProofOfWork(b *Block) *ProofOfWork {
  target := big.NewInt(1)
  target.Lsh(target, uint(256 - targetBits))
  return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) PrepareData(nonce int64) []byte {
  data := bytes.Join([][]byte{
      pow.block.PrevBlockHash,
      pow.block.Data,
      IntToHex(pow.block.Timestamp),
      IntToHex(int64(targetBits)),
      IntToHex(int64(nonce)),
    },[]byte{})
  return data
}

func (pow *ProofOfWork) Mine() (int64, []byte) {
  var hashInt big.Int
  var hash [32]byte
  nonce := int64(0)
  fmt.Printf("Mining the block containing \"%s\"\n",pow.block.Data)
  for nonce < maxNonce {
    data := pow.PrepareData(nonce)
    hash = sha256.Sum256(data)
    fmt.Printf("\r%x",hash)
    hashInt.SetBytes(hash[:])
    if(hashInt.Cmp(pow.target) == -1){
      break
    } else {
      nonce ++
    }
  }
  fmt.Printf("\n\n")
  return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
  var hashInt big.Int
  data := pow.PrepareData(pow.block.Nonce)
  hash := sha256.Sum256(data)
  hashInt.SetBytes(hash[:])
  isValid := hashInt.Cmp(pow.target) == -1
  return isValid
}
