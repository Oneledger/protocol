package main

import (
	"context"
	"contracts/testgas"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	//"contracts/erc20" // for demo
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	//ctx := context.Background()
	privateKey, err := crypto.HexToECDSA("6ba785c66db47c5d9266ca4c351c90e238300f39122e7d30a3148f1bd95f0c5b")
	if err != nil {
		log.Fatal(err)
	}



	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}


	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}


	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("suggest gas price")
		log.Fatal(err)
	}



	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(6721970) // in units
	auth.GasPrice = gasPrice

	//tokenaddress := common.HexToAddress("0x01Ff2f46A168A35072Df78129cB64D53af324641")
	//oltaddress := common.HexToAddress("0x7488974D8469234f1eb251C2aa05D4403543Df13")  // also used to deploy contract
	validatorAddress := common.HexToAddress("0xEF66117e2014cEb509A039b5ADeBed7012332935")
	//nameSource := []byte("ERC20Basic")
	//var name [32]byte
	//copy(nameSource[:], name[:32])

	address, tx, instance, err := testgas.DeployTestRedeemGas(auth, client, []common.Address{validatorAddress})
	//var x, _ = new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	//address, tx, instance, err := erc20.DeployERC20Basic(auth, client,x)
	if err != nil {
		fmt.Println("deploy contract error")
		log.Fatal(err)
	}

	session := testgas.TestRedeemGasSession{
		Contract: instance,
		CallOpts: bind.CallOpts{
			Pending: true,        // Acts on pending state if set to true
			From: auth.From,
			Context: context.Background(),
		},
		TransactOpts: bind.TransactOpts{
			From: auth.From,
			Nonce: nil,           // nil uses nonce of pending state
			Signer: auth.Signer,
			Value: big.NewInt(0),
			GasPrice: nil,        // nil automatically suggests gas price
			GasLimit: 0,          // 0 automatically estimates gas limit
		},
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	valid, err := session.IsValidator(validatorAddress)
	if err != nil {
		log.Printf("Error in qureying deployed contract: %v\n", err)
		return
	}
	fmt.Printf("Is Validater: %t\n", valid)
	//message := "TestRedeem"
	//hashRaw := crypto.Keccak256([]byte(message))
	//signature, err := crypto.Sign(hashRaw, privateKey)
    //stringSig := hex.EncodeToString(signature)
    //fmt.Println(stringSig)
    //s := base64.StdEncoding.Encode(stigSig2, signature)
	//validationMsg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message
	//messagehash := hex.EncodeToString(hashRaw)
    //r := stringToKeccak256(stringSig[:32])
    //s := stringToKeccak256(stringSig[32:64])
    //v := stringToKeccak256(string(int(stringSig[65]) + 27) )
    //fmt.Println(validatorPayload)
	signTX,err := session.Sign()
	if err != nil {
		log.Printf("Error in qureying deployed contract: %v\n", err)
		return
	}
	fmt.Printf("Validtor Signed ,Wait for trasaction hash to confirm :%s \n", signTX.Hash().Hex() )


    //validatorPayload := [stringToKeccak256("test"),stringToKeccak256("00")]
    //v,err := session.OfChainSign()

	//ERC20 token
    //Adress : 0x01Ff2f46A168A35072Df78129cB64D53af324641
	//TXHASH : 0x89869d527033c1dce14bcc81a753486565c077267e6048e946bc1b4fed6571e1

	//LOCKREDEEM
	//Address : 0xACDD9ad30b4Dc12e2b595888c2174B84E5E626Cd
	//TXHASH : 0x7787b04ce35db29a354d00192926f2f6172433443b0211a19d2595815eb298d4


	_ = instance
}
