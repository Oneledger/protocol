package htlc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"time"

	"github.com/Oneledger/protocol/node/chains/bitcoin/rpc"
	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"golang.org/x/crypto/ripemd160"
)

const verify = true

const secretSize = 32

const txVersion = 2

// There are two directions that the atomic swap can be performed, as the
// initiator can be on either chain.  This tool only deals with creating the
// Bitcoin transactions for these swaps.  A second tool should be used for the
// transaction on the other chain.  Any chain can be used so long as it supports
// OP_SHA256 and OP_CHECKLOCKTIMEVERIFY.
//
// Example scenerios using bitcoin as the second chain:
//
// Scenerio 1:
//   cp1 initiates (olt)
//   cp2 participates with cp1 H(S) (btc)
//   cp1 redeems btc revealing S
//     - must verify H(S) in contract is hash of known Secret
//   cp2 redeems olt with S
//
// Scenerio 2:
//   cp1 initiates (btc)
//   cp2 participates with cp1 H(S) (olt)
//   cp1 redeems olt revealing S
//     - must verify H(S) in contract is hash of known Secret
//   cp2 redeems btc with S

//type Command interface {
//	RunCommand(*rpc.Bitcoind) error
//}
//
//// offline commands don't require wallet RPC.
//type OfflineCommand interface {
//	Command
//	RunOfflineCommand() error
//}

type InitiateCmd struct {
	cp2Addr  	*btcutil.AddressPubKeyHash
	amount   	btcutil.Amount
	lockTime 	int64
	scrHash  	[secretSize]byte
	Contract 	[]byte
	ContractTx	*wire.MsgTx
	RefundTx	*wire.MsgTx
}

type RedeemCmd struct {
	contract   			[]byte
	contractTx 			*wire.MsgTx
	secret     			[]byte
	RedeemContractTx	*wire.MsgTx
}

type ParticipateCmd struct {
	cp1Addr    *btcutil.AddressPubKeyHash
	amount     btcutil.Amount
	secretHash []byte
	lockTime   int64
	Contract 	[]byte
	ContractTx	*wire.MsgTx
}

type RefundCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

type ExtractSecretCmd struct {
	redemptionTx *wire.MsgTx
	secretHash   []byte
	Secret		 []byte
}

type AuditContractCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
	SecretHash [secretSize]byte
}

func normalizeAddress(addr string, defaultPort string) (hostport string, err error) {
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err = net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}

func walletPort(params *chaincfg.Params) string {
	switch params {
	case &chaincfg.MainNetParams:
		return "8332"
	case &chaincfg.TestNet3Params:
		return "18332"
	default:
		return ""
	}
}

// contractArgs specifies the common parameters used to create the initiator's
// and participant's contract.
type contractArgs struct {
	them       *btcutil.AddressPubKeyHash
	amount     btcutil.Amount
	locktime   int64
	secretHash []byte
}

// atomicSwapContract returns an output script that may be redeemed by one of
// two signature scripts:
//
//   <their sig> <their pubkey> <initiator Secret> 1
//
//   <my sig> <my pubkey> 0
//
// The first signature script is the normal redemption path done by the other
// party and requires the initiator's Secret.  The second signature script is
// the refund path performed by us, but the refund can only be performed after
// locktime.
func atomicSwapContract(pkhMe, pkhThem *[ripemd160.Size]byte, locktime int64, secretHash []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()

	b.AddOp(txscript.OP_IF) // Normal redeem path [0]
	{
		// Require initiator's Secret to be a known length that the redeeming
		// party can audit.  This is used to prevent fraud attacks between two
		// currencies that have different maximum data sizes.
		b.AddOp(txscript.OP_SIZE)  					//[1]
		b.AddInt64(secretSize)						//[2]
		b.AddOp(txscript.OP_EQUALVERIFY)			//[3]

		// Require initiator's Secret to be known to redeem the output.
		b.AddOp(txscript.OP_SHA256)					//[4]
		b.AddData(secretHash)						//[5]
		b.AddOp(txscript.OP_EQUALVERIFY)			//[6]

		// Verify their signature is being used to redeem the output.  This
		// would normally end with OP_EQUALVERIFY OP_CHECKSIG but this has been
		// moved outside of the branch to save a couple bytes.
		b.AddOp(txscript.OP_DUP)					//[7]
		b.AddOp(txscript.OP_HASH160)				//[8]
		b.AddData(pkhThem[:])						//[9]
	}
	b.AddOp(txscript.OP_ELSE) // Refund path		//[10]
	{
		// Verify locktime and drop it off the stack (which is not done by
		// CLTV).
		b.AddInt64(locktime)						//[11]
		b.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)	//[12]
		b.AddOp(txscript.OP_DROP)					//[13]

		// Verify our signature is being used to refund the output.  This would
		// normally end with OP_EQUALVERIFY OP_CHECKSIG but this has been moved
		// outside of the branch to save a couple bytes.
		b.AddOp(txscript.OP_DUP)					//[14]
		b.AddOp(txscript.OP_HASH160)				//[15]
		b.AddData(pkhMe[:])							//[16]
	}
	b.AddOp(txscript.OP_ENDIF)						//[17]

	// Complete the signature check.
	b.AddOp(txscript.OP_EQUALVERIFY)				//[18]
	b.AddOp(txscript.OP_CHECKSIG)					//[19]

	return b.Script()
}

// builtContract houses the details regarding a contract and the contract
// payment transaction, as well as the transaction to perform a refund.
type builtContract struct {
	contract       []byte
	contractP2SH   btcutil.Address
	contractTxHash *chainhash.Hash
	contractTx     *wire.MsgTx
	contractFee    btcutil.Amount
	refundTx       *wire.MsgTx
	refundFee      btcutil.Amount
}

// buildContract creates a contract for the parameters specified in args, using
// wallet RPC to generate an internal address to redeem the refund and to sign
// the payment to the contract transaction.
func buildContract(b *rpc.Bitcoind, args *contractArgs) (*builtContract, error) {
	refundAddr, err := b.GetRawChangeAddress()
	if err != nil {
		return nil, fmt.Errorf("getrawchangeaddress: %v", err)
	}
	refundAddrH, ok := refundAddr.(interface{ Hash160() *[ripemd160.Size]byte })
	if !ok {
		return nil, errors.New("unable to create hash160 from change address")
	}

	contract, err := atomicSwapContract(refundAddrH.Hash160(), args.them.Hash160(),
		args.locktime, args.secretHash)
	if err != nil {
		return nil, err
	}
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, b.ChainParams)
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
	}

	feePerKb, minFeePerKb, err := b.GetFeePerKb()
	if err != nil {
		return nil, err
	}

	unsignedContract := wire.NewMsgTx(txVersion)
	unsignedContract.AddTxOut(wire.NewTxOut(int64(args.amount), contractP2SHPkScript))
	unsignedContract, contractFee, err := b.FundRawTransaction(unsignedContract, feePerKb)
	if err != nil {
		return nil, fmt.Errorf("fundrawtransaction: %v", err)
	}
	contractTx, complete, err := b.SignRawTransaction(unsignedContract)
	if err != nil {
		return nil, fmt.Errorf("signrawtransaction: %v", err)
	}
	if !complete {
		return nil, errors.New("signrawtransaction: failed to completely sign contract transaction")
	}

	contractTxHash := contractTx.TxHash()

	refundTx, refundFee, err := buildRefund(b, contract, contractTx, feePerKb, minFeePerKb)
	if err != nil {
		return nil, err
	}

	return &builtContract{
		contract,
		contractP2SH,
		&contractTxHash,
		contractTx,
		contractFee,
		refundTx,
		refundFee,
	}, nil
}

// createSig creates and returns the serialized raw signature and compressed
// pubkey for a transaction input signature.  Due to limitations of the Bitcoin
// Core RPC API, this requires dumping a private key and signing in the client,
// rather than letting the wallet sign.

func createSig(b *rpc.Bitcoind, tx *wire.MsgTx, idx int, pkScript []byte, address btcutil.Address) (sig, pubkey []byte, err error) {
	wif, err := b.DumpPrivKey(address)
	if err != nil {
		return nil, nil, err
	}
	sig, err = txscript.RawTxInSignature(tx, idx, pkScript, txscript.SigHashAll, wif.PrivKey)
	if err != nil {
		return nil, nil, err
	}
	return sig, wif.PrivKey.PubKey().SerializeCompressed(), nil
}

func buildRefund(b *rpc.Bitcoind, contract []byte, contractTx *wire.MsgTx, feePerKb, minFeePerKb btcutil.Amount) (

	refundTx *wire.MsgTx, refundFee btcutil.Amount, err error) {

	contractP2SH, err := btcutil.NewAddressScriptHash(contract, b.ChainParams)
	if err != nil {
		return nil, 0, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, 0, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{Hash: contractTxHash, Index: ^uint32(0)}
	for i, o := range contractTx.TxOut {
		if bytes.Equal(o.PkScript, contractP2SHPkScript) {
			contractOutPoint.Index = uint32(i)
			break
		}
	}
	if contractOutPoint.Index == ^uint32(0) {
		return nil, 0, errors.New("contract tx does not contain a P2SH contract payment")
	}

	refundAddress, err := b.GetRawChangeAddress()
	if err != nil {
		return nil, 0, fmt.Errorf("getrawchangeaddress: %v", err)
	}

	refundOutScript, err := txscript.PayToAddrScript(refundAddress)
	if err != nil {
		return nil, 0, err
	}

	log.Debug("About to Extract", "contract", hex.EncodeToString(contract))
	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		log.Fatal("ExtractAtomicSwapDataPushes", "err", err)
	}
	if pushes == nil {
		log.Warn("Not a swap contract")
	}

	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], b.ChainParams)
	if err != nil {
		return nil, 0, err
	}

	refundTx = wire.NewMsgTx(txVersion)
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript)) // amount set below

	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee = txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutPoint.Index].Value - int64(refundFee)

	if txrules.IsDustOutput(refundTx.TxOut[0], minFeePerKb) {
		return nil, 0, fmt.Errorf("refund output value of %v is dust", btcutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutPoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubKey, err := createSig(b, refundTx, 0, contract, refundAddr)
	log.Debug("createsig refund", "refundaddr", refundAddr, "contract", contract)
	if err != nil {
		return nil, 0, err
	}

	log.Debug("About to Refund")
	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubKey)
	if err != nil {
		return nil, 0, err
	}
	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), contractTx.TxOut[contractOutPoint.Index].Value)
		if err != nil {
			panic(err)
		}
		err = e.Execute()
		if err != nil {
			panic(err)
		}
	}

	log.Debug("Finished")
	return refundTx, refundFee, nil
}

func sha256Hash(x []byte) []byte {
	h := sha256.Sum256(x)
	return h[:]
}

func calcFeePerKb(absoluteFee btcutil.Amount, serializeSize int) float64 {
	return float64(absoluteFee) / float64(serializeSize) / 1e5
}

func copyArray(from []byte) []byte {
	to := make([]byte, len(from))
	for i := 0; i < len(from); i++ {
		to[i] = from[i]
	}
	return to
}

func copyMsgTx(from *wire.MsgTx) *wire.MsgTx {
	var to wire.MsgTx

	to.Version = from.Version
	to.TxIn = copyTxIn(from.TxIn)
	to.TxOut = copyTxOut(from.TxOut)
	to.LockTime = from.LockTime

	return &to
}

func copyTxIn(from []*wire.TxIn) []*wire.TxIn {
	var to []*wire.TxIn
	to = make([]*wire.TxIn, len(from))
	for i := 0; i < len(from); i++ {
		to[i] = &wire.TxIn{}
		to[i].PreviousOutPoint = from[i].PreviousOutPoint
		to[i].SignatureScript = copyArray(from[i].SignatureScript)
		to[i].Witness = from[i].Witness
		to[i].Sequence = from[i].Sequence
	}
	return to
}

func copyTxOut(from []*wire.TxOut) []*wire.TxOut {
	var to []*wire.TxOut
	to = make([]*wire.TxOut, len(from))
	for i := 0; i < len(from); i++ {
		to[i] = &wire.TxOut{}
		to[i].Value = from[i].Value
		to[i].PkScript = copyArray(from[i].PkScript)
	}
	return to
}

func (cmd *InitiateCmd) RunCommand(c *rpc.Bitcoind) (*chainhash.Hash, error) {
	log.Debug("About to Initiate")
	//var Secret [secretSize]byte
	//_, err := rand.Read(Secret[:])
	//if err != nil {
	//	return nil, err
	//}
	secretHash := cmd.scrHash[:]

	log.Debug("About to build contract")
	b, err := buildContract(c, &contractArgs{
		them:       cmd.cp2Addr,
		amount:     cmd.amount,
		locktime:   cmd.lockTime,
		secretHash: secretHash,
	})
	if err != nil {
		return nil, err
	}
	cmd.Contract = copyArray(b.contract)
	cmd.ContractTx = copyMsgTx(b.contractTx)
	cmd.RefundTx = copyMsgTx(b.refundTx)

	refundTxHash := b.refundTx.TxHash()
	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())
	refundFeePerKb := calcFeePerKb(b.refundFee, b.refundTx.SerializeSize())

	log.Debug("About to grow")

	fmt.Printf("Secret:      %x\n", "unknown")
	fmt.Printf("Secret hash: %x\n\n", secretHash)
	fmt.Printf("Contract fee: %v (%0.8f BTC/kB)\n", b.contractFee, contractFeePerKb)
	fmt.Printf("Refund fee:   %v (%0.8f BTC/kB)\n\n", b.refundFee, refundFeePerKb)
	fmt.Printf("Contract (%v):\n", b.contractP2SH)
	fmt.Printf("%x\n\n", b.contract)
	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)
	fmt.Printf("Contract transaction (%v):\n", b.contractTxHash)
	fmt.Printf("%x\n\n", contractBuf.Bytes())
	var refundBuf bytes.Buffer
	refundBuf.Grow(b.refundTx.SerializeSize())
	b.refundTx.Serialize(&refundBuf)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", refundBuf.Bytes())

	log.Debug("About to Publish")
	return c.PublishTx(b.contractTx, "contract")
}

// refundP2SHContract returns the signature script to refund a contract output
// using the contract author's signature after the locktime has been reached.
// This function assumes P2SH and appends the contract as the final data push.
func refundP2SHContract(contract, sig, pubkey []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddInt64(0)
	b.AddData(contract)
	return b.Script()
}

func (cmd *RedeemCmd) RunCommand(c *rpc.Bitcoind) (*chainhash.Hash, error) {
	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		log.Debug("Extract")
		return nil, err
	}
	if pushes == nil {
		return nil, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		c.ChainParams)
	if err != nil {
		return nil, err
	}
	contractHash := btcutil.Hash160(cmd.contract)
	contractOut := -1
	for i, out := range cmd.contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, c.ChainParams)
		if sc == txscript.ScriptHashTy &&
			bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return nil, errors.New("transaction does not contain a contract output")
	}

	addr, err := c.GetRawChangeAddress()
	if err != nil {
		log.Debug("RawChange")
		return nil, err
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		log.Debug("PayToAddr")
		return nil, err
	}

	contractTxHash := cmd.contractTx.TxHash()
	contractOutPoint := wire.OutPoint{
		Hash:  contractTxHash,
		Index: uint32(contractOut),
	}

	feePerKb, minFeePerKb, err := c.GetFeePerKb()
	if err != nil {
		log.Debug("GetFee")
		return nil, err
	}

	redeemTx := wire.NewMsgTx(txVersion)
	redeemTx.LockTime = uint32(pushes.LockTime)
	redeemTx.AddTxIn(wire.NewTxIn(&contractOutPoint, nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, outScript)) // amount set below
	redeemSize := estimateRedeemSerializeSize(cmd.contract, redeemTx.TxOut)
	fee := txrules.FeeForSerializeSize(feePerKb, redeemSize)
	redeemTx.TxOut[0].Value = cmd.contractTx.TxOut[contractOut].Value - int64(fee)
	if txrules.IsDustOutput(redeemTx.TxOut[0], minFeePerKb) {
		return nil, fmt.Errorf("redeem output value of %v is dust", btcutil.Amount(redeemTx.TxOut[0].Value))
	}

	redeemSig, redeemPubKey, err := createSig(c, redeemTx, 0, cmd.contract, recipientAddr)
	if err != nil {
		log.Debug("createSig")
		return nil, err
	}
	redeemSigScript, err := redeemP2SHContract(cmd.contract, redeemSig, redeemPubKey, cmd.secret)
	if err != nil {
		log.Debug("redeem", "err", err)
		return nil, err
	}
	redeemTx.TxIn[0].SignatureScript = redeemSigScript

	redeemTxHash := redeemTx.TxHash()
	redeemFeePerKb := calcFeePerKb(fee, redeemTx.SerializeSize())

	var buf bytes.Buffer
	buf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&buf)
	fmt.Printf("Redeem fee: %v (%0.8f BTC/kB)\n\n", fee, redeemFeePerKb)
	fmt.Printf("Redeem transaction (%v):\n", &redeemTxHash)
	fmt.Printf("%x\n\n", buf.Bytes())

	cmd.RedeemContractTx = copyMsgTx(redeemTx)

	if verify {
		e, err := txscript.NewEngine(cmd.contractTx.TxOut[contractOutPoint.Index].PkScript,
			redeemTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), cmd.contractTx.TxOut[contractOut].Value)
		if err != nil {
			log.Debug("NewEngine", "err", err)
			panic(err)
		}
		err = e.Execute()
		if err != nil {
			log.Debug("Execute", "err", err)
			panic(err)
		}
	}

	return c.PublishTx(redeemTx, "redeem")
}

// redeemP2SHContract returns the signature script to redeem a contract output
// using the redeemer's signature and the initiator's Secret.  This function
// assumes P2SH and appends the contract as the final data push.
func redeemP2SHContract(contract, sig, pubkey, secret []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddData(secret)
	b.AddInt64(1)
	b.AddData(contract)
	return b.Script()
}

func (cmd *ParticipateCmd) RunCommand(c *rpc.Bitcoind) (*chainhash.Hash, error) {
	b, err := buildContract(c, &contractArgs{
		them:       cmd.cp1Addr,
		amount:     cmd.amount,
		locktime:   cmd.lockTime,
		secretHash: cmd.secretHash,
	})
	if err != nil {
		return nil, err
	}

	cmd.Contract = copyArray(b.contract)
	cmd.ContractTx = copyMsgTx(b.contractTx)

	refundTxHash := b.refundTx.TxHash()
	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())
	refundFeePerKb := calcFeePerKb(b.refundFee, b.refundTx.SerializeSize())

	fmt.Printf("Contract fee: %v (%0.8f BTC/kB)\n", b.contractFee, contractFeePerKb)
	fmt.Printf("Refund fee:   %v (%0.8f BTC/kB)\n\n", b.refundFee, refundFeePerKb)
	fmt.Printf("Contract (%v):\n", b.contractP2SH)
	fmt.Printf("%x\n\n", b.contract)
	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)
	fmt.Printf("Contract transaction (%v):\n", b.contractTxHash)
	fmt.Printf("%x\n\n", contractBuf.Bytes())
	var refundBuf bytes.Buffer
	refundBuf.Grow(b.refundTx.SerializeSize())
	b.refundTx.Serialize(&refundBuf)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", refundBuf.Bytes())

	return c.PublishTx(b.contractTx, "contract")
}

func (cmd *RefundCmd) RunCommand(c *rpc.Bitcoind) (*chainhash.Hash, error) {
	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return nil, err
	}
	if pushes == nil {
		return nil, errors.New("contract is not an atomic swap script recognized by this tool")
	}

	feePerKb, minFeePerKb, err := c.GetFeePerKb()
	if err != nil {
		return nil, err
	}

	refundTx, refundFee, err := buildRefund(c, cmd.contract, cmd.contractTx, feePerKb, minFeePerKb)
	if err != nil {
		return nil, err
	}
	refundTxHash := refundTx.TxHash()
	var buf bytes.Buffer
	buf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&buf)

	refundFeePerKb := calcFeePerKb(refundFee, refundTx.SerializeSize())

	fmt.Printf("Refund fee: %v (%0.8f BTC/kB)\n\n", refundFee, refundFeePerKb)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", buf.Bytes())

	return c.PublishTx(refundTx, "refund")
}

func (cmd *ExtractSecretCmd) RunCommand(c *rpc.Bitcoind) error {
	return cmd.RunOfflineCommand()
}

func (cmd *ExtractSecretCmd) RunOfflineCommand() error {
	// Loop over all pushed data from all inputs, searching for one that hashes
	// to the expected hash.  By searching through all data pushes, we avoid any
	// issues that could be caused by the initiator redeeming the participant's
	// contract with some "nonstandard" or unrecognized transaction or script
	// type.
	for _, in := range cmd.redemptionTx.TxIn {
		pushes, err := txscript.PushedData(in.SignatureScript)
		if err != nil {
			return err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), cmd.secretHash) {
				fmt.Printf("Secret: %x\n", push)
				cmd.Secret = copyArray(push)
				return nil
			}
		}
	}
	return errors.New("transaction does not contain the Secret")
}

func (cmd *AuditContractCmd) RunCommand(b *rpc.Bitcoind) error {
	log.Debug("======================================================", "bitcoin", b)
	log.Debug("Audit (Offline)", "bitcoin", b)
	return cmd.RunOfflineCommand(b.ChainParams)
}

func (cmd *AuditContractCmd) RunOfflineCommand(chainParams *chaincfg.Params) error {
	contractHash160 := btcutil.Hash160(cmd.contract)
	contractOut := -1
	for i, out := range cmd.contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool")
	}
	if pushes.SecretSize != secretSize {
		return fmt.Errorf("contract specifies strange Secret size %v", pushes.SecretSize)
	}

	contractAddr, err := btcutil.NewAddressScriptHash(cmd.contract, chainParams)
	if err != nil {
		return err
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return err
	}
	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		chainParams)
	if err != nil {
		return err
	}

	fmt.Printf("Contract address:        %v\n", contractAddr)
	fmt.Printf("Contract value:          %v\n", btcutil.Amount(cmd.contractTx.TxOut[contractOut].Value))
	fmt.Printf("Recipient address:       %v\n", recipientAddr)
	fmt.Printf("Author's refund address: %v\n\n", refundAddr)

	fmt.Printf("Secret hash: %x\n\n", pushes.SecretHash[:])
	cmd.SecretHash = pushes.SecretHash
	if pushes.LockTime >= int64(txscript.LockTimeThreshold) {
		t := time.Unix(pushes.LockTime, 0)
		fmt.Printf("Locktime: %v\n", t.UTC())
		reachedAt := time.Until(t).Truncate(time.Second)
		if reachedAt > 0 {
			fmt.Printf("Locktime reached in %v\n", reachedAt)
		} else {
			fmt.Printf("Contract refund time lock has expired\n")
		}
	} else {
		fmt.Printf("Locktime: block %v\n", pushes.LockTime)
	}

	return nil
}
