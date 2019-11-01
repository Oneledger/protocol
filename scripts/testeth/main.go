package main

import (
	"fmt"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"os"

	//"github.com/Oneledger/protocol/log"
	//"os"
)

func main() {
	var initialValidators = []string{"0xa5d180c3be91e70cb00ca3a2b67fe2664ae61087"}
	var keypath = "/scripts/testeth/keydump"
	var contractaddress = "0xc8fb8b04747c392ca84da4f8ac26521c3cef4f9e"
	config := ethereum.CreateEthConfig("http://127.0.0.1:8545", keypath,contractaddress,initialValidators)
	var log = log.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	access,err := ethereum.NewAccess("/home/tanmay/Codebase/protocol",config,log)

	if(err!=nil){
		fmt.Println(err)
	}
	fmt.Println(access)
	balance,err := access.Balance()
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
	var weiamout = new(big.Int)
	weiamout.SetString("10000000000000000000", 10)
	locktx,err := access.Lock(weiamout)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(locktx)
	toAddress := common.HexToAddress("0x313Df6717f5D9B9d387758B0B1e32928B6BeAf51")
	signtx,err := access.Sign(weiamout, toAddress)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(signtx)


}
