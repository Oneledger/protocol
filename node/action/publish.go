/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, publish, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
    "bytes"
    "github.com/Oneledger/protocol/node/data"
    "github.com/Oneledger/protocol/node/chains/bitcoin"
    "github.com/Oneledger/protocol/node/chains/ethereum"
    "github.com/Oneledger/protocol/node/comm"

)

// Synchronize a publish between two users
type Publish struct {
	Base

	Target     id.AccountKey `json:"target"`
	Contract   Message       `json:"message"` //message converted from HTLContract
	SecretHash [32]byte      `json:"secrethash"`
	Sequence   int           `json:"sequence"`
}

// Ensure that all of the base values are at least reasonable.
func (transaction *Publish) Validate() err.Code {
	log.Debug("Validating Publish Transaction")

	if transaction.Target == nil {
		log.Debug("Missing Target")
		return err.MISSING_DATA
	}

	if transaction.Contract == nil {
		log.Debug("Missing Contract")
		return err.MISSING_DATA
	}

	log.Debug("Publish is validated!")
	return err.SUCCESS
}

func (transaction *Publish) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Publish Transaction for CheckTx")

	// TODO: Check all of the data to make sure it is valid.

	return err.SUCCESS
}

// Start the publish
func (transaction *Publish) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Publish Transaction for DeliverTx")

    commands := transaction.Expand(app)

    transaction.Resolve(app, commands)

    //before loop of execute, lastResult is nil
    var lastResult map[Parameter]FunctionValue

    for i := 0; i < commands.Count(); i++ {
        status, result := Execute(app, commands[i], lastResult)
        if status != err.SUCCESS {
            log.Error("Failed to Execute", "command", commands[i])
            return err.EXPAND_ERROR
        }
        lastResult = result
    }
	return err.SUCCESS
}

// Is this node one of the partipants in the publish
func (transaction *Publish) ShouldProcess(app interface{}) bool {
	account := transaction.GetNodeAccount(app)

	if bytes.Equal(transaction.Target, account.AccountKey()) {
		return true
	}

	return false
}

func (transaction *Publish) GetNodeAccount(app interface{}) id.Account {

	accounts := GetAccounts(app)
	account, _ := accounts.FindName(global.Current.NodeAccountName)
	if account == nil {
		log.Error("Node does not have account", "name", global.Current.NodeAccountName)
		accounts.Dump()
		return nil
	}

	return account
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Publish) Expand(app interface{}) Commands {
    swap := transaction.FindSwap(app)
    account := transaction.GetNodeAccount(app)
    isParty := swap.IsParty(account)
    role := swap.getRole(*isParty)
    chains := swap.getChains()
    if transaction.Sequence > 1 {
        role = ALL
        log.Debug("Publish role", "role", role)
    }
	return GetCommands(PUBLISH, role, chains)
}

func (transaction *Publish) FindSwap(app interface{}) *Swap {

    status := GetStatus(app)
    senderKey := transaction.Base.Owner
    log.Debug("FindSwap", "key",senderKey)
    swap := FindSwap(status, senderKey).(*Swap)
    return swap
}

// Plug in data from the rest of a system into a set of commands
func (publish *Publish) Resolve(app interface{}, commands Commands) {
	swap := publish.FindSwap(app)
	swap.Resolve(app, commands)

	for i := 0; i < len(commands); i++ {

	    if commands[i].Function == AUDITCONTRACT || commands[i].Function == EXTRACTSECRET {
            if commands[i].Chain == data.BITCOIN {
                contract := bitcoin.GetHTLCFromMessage(publish.Contract)
                commands[i].Data[BTCCONTRACT] = contract

            } else if commands[i].Chain == data.ETHEREUM {
                contract := ethereum.GetHTLCFromMessage(publish.Contract)
                commands[i].Data[ETHCONTRACT] = contract
            }
        }
        commands[i].Data[PREIMAGE] = publish.SecretHash
        if commands[i].Function == SUBMIT_TRANSACTION {
            commands[i].Data[SEQUENCE] = publish.Sequence + 1
            commands[i].Data[CHAINID] = GetChainID(app)
        }
        //log.Debug("resolved command", "command", commands[i], "sequence", commands[i].Data[SEQUENCE])
    }
    return
}


func SubmitTransactionOLT(context map[Parameter]FunctionValue, chain data.ChainType) (bool, map[Parameter]FunctionValue) {
    signers := make([]PublicKey, 0)
    owner := GetParty(context[MY_ACCOUNT])
    target := GetParty(context[THEM_ACCOUNT])

    var contract Message
    if chain == data.BITCOIN {
        contract = GetBTCContract(context[BTCCONTRACT]).ToMessage()

    } else if chain == data.ETHEREUM {
        contract = GetETHContract(context[ETHCONTRACT]).ToMessage()
    }

    sequence := GetInt(context[SEQUENCE])
    secretHash := GetByte32(context[PREIMAGE])
    chainId := GetString(context[CHAINID])
    //log.Debug("parsed contract", "contract", contract, "chain", chain, "context", context, "sequence", sequence)
    publish := &Publish{
       Base: Base{
           Type:     PUBLISH,
           ChainId:  chainId,
           Signers:  signers,
           Owner:    owner.Key,
           Sequence: global.Current.Sequence,
       },
       Target:     target.Key,
       Contract:   contract,
       SecretHash: secretHash,
       Sequence:   sequence,
    }

    packet := SignAndPack(PUBLISH, Transaction(publish))

    result := comm.Broadcast(packet)
    log.Debug("Submit Transaction to OLT successfully", "result", result)
    return true, nil
}