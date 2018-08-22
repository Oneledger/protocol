package ethereum

import (
    "testing"
    "github.com/Oneledger/protocol/node/chains/ethereum/htlc"
    "github.com/ethereum/go-ethereum/common"
    "github.com/Oneledger/protocol/node/log"

    "context"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/ethclient"
    "strings"
)

func TestHTLContract_Refund(t *testing.T) {

    address := common.BytesToAddress([]byte("0xe86778b054193b21186d9fb3ffedf1ce042e347d"))

    key := `{"address":"aafa2d8980a730b02195f9c8dfeafeb3e69a69ca","crypto":{"cipher":"aes-128-ctr","ciphertext":"cfc36b7deb503116482371b7d2596aa936758b8247279efce461cf0344ae4b31","cipherparams":{"iv":"fc200b937116258856dd0e5a085e011d"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"7354c4523dfc70372c8c34616c15dc21448ac40617ffc3a7b3a9af7ee32c37e6"},"mac":"1b322fa3c5789cede87144783f2bd8c4588e5094e68f7880640ca9a5458b8aab"},"id":"87363a39-0171-4640-ba12-b5aacad7aed2","version":3}`
    auth, err := bind.NewTransactor(strings.NewReader(key), "2345")

    cli, err := ethclient.Dial("/home/lan/go/test/ethereum/B/geth.ipc")
    if err != nil {
        log.Error("failed to get geth ipc ", "err", err)
    }
    contract, err := htlc.NewHtlc(address, cli)
    if err != nil {
        log.Error("can't get new htlc", "err", err)
        return
    }
    tx, err := contract.Refund(auth)
    if err != nil {
        log.Error("refund failed", "err", err)
        return
    }

    ctx := context.Background()
    receipt, err := cli.TransactionReceipt(ctx, tx.Hash())
    if err != nil {
        log.Error("Failed to get the receipt", "err", err)
        return
    }
    if receipt.Status == 0 {
        log.Error("setup failure","status", receipt.Status)
        return
    }
    balance, _ := contract.Balance(&bind.CallOpts{Pending: false})
    log.Info("balance", "balance", balance)

}
