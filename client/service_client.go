package client

import (
	"errors"
	"os"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/rpc"
)

// A type-safe client for accessing rpc services.
// Eventually each service will be broken out onto its own type
// TODO: The methods defined here should handle context.Context
type ServiceClient struct {
	*rpc.Client
}

func NewServiceClient(conn string) (*ServiceClient, error) {
	client, err := rpc.NewClient(conn)
	return &ServiceClient{client}, err
}

func (c *ServiceClient) Balance(addr keys.Address) (out BalanceReply, err error) {
	/*if len(request) <= 20 {
		return out, errors.New("address has insufficient length")
	}*/
	request := BalanceRequest{addr}
	err = c.Call("query.Balance", &request, &out)
	return
}

func (c *ServiceClient) BalancePool(poolname string) (out BalanceReply, err error) {
	request := BalancePoolRequest{poolname}
	err = c.Call("query.BalancePool", &request, &out)
	return
}

func (c *ServiceClient) ValidatorStatus(request ValidatorStatusRequest) (out ValidatorStatusReply, err error) {
	err = c.Call("query.ValidatorStatus", &request, &out)
	return
}

func (c *ServiceClient) DelegationStatus(request DelegationStatusRequest) (out DelegationStatusReply, err error) {
	err = c.Call("query.DelegationStatus", &request, &out)
	return
}

func (c *ServiceClient) EVMCall(request SendTxRequest) (out EVMCallReply, err error) {
	err = c.Call("query.EVMCall", &request, &out)
	return
}

func (c *ServiceClient) EVMAccount(request EVMAccountRequest) (out EVMAccountReply, err error) {
	err = c.Call("query.EVMAccount", &request, &out)
	return
}

func (c *ServiceClient) EVMEstimateGas(request SendTxRequest) (out EVMEstimateGasReply, err error) {
	err = c.Call("query.EVMEstimateGas", &request, &out)
	return
}

func (c *ServiceClient) EVMTransactionLogs(request EVMTransactionLogsRequest) (out EVMLogsReply, err error) {
	err = c.Call("query.EVMTransactionLogs", &request, &out)
	return
}

func (c *ServiceClient) CurrBalance(addr keys.Address, currency string) (out CurrencyBalanceReply, err error) {
	/*if len(request) <= 20 {
		return out, errors.New("address has insufficient length")
	}*/
	request := CurrencyBalanceRequest{currency, addr}
	err = c.Call("query.CurrencyBalance", &request, &out)
	return
}

func (c *ServiceClient) NodeName() (out NodeNameReply, err error) {
	err = c.Call("node.NodeName", struct{}{}, &out)
	return
}

func (c *ServiceClient) NodeAddress() (out NodeAddressReply, err error) {
	err = c.Call("node.Address", struct{}{}, &out)
	return
}

func (c *ServiceClient) NodeID(req NodeIDRequest) (out NodeIDReply, err error) {
	err = c.Call("node.ID", req, &out)
	return
}

func (c *ServiceClient) SendTx(req SendTxRequest) (out CreateTxReply, err error) {
	if os.Getenv("OLTEST") == "1" {
		err = c.Call("tx.SendTx", req, &out)
	} else {
		err = errors.New("SendTx disabled")
	}
	return
}

func (c *ServiceClient) AddAccount(req AddAccountRequest) (out AddAccountReply, err error) {
	err = c.Call("owner.AddAccount", req, &out)
	return
}

func (c *ServiceClient) DeleteAccount(req DeleteAccountRequest) (out DeleteAccountReply, err error) {
	err = c.Call("owner.DeleteAccount", req, &out)
	return
}

func (c *ServiceClient) ListAccounts() (out ListAccountsReply, err error) {
	err = c.Call("owner.ListAccounts", struct{}{}, &out)
	return
}

func (c *ServiceClient) ListAccountAddresses() (out ListAccountAddressesReply, err error) {
	err = c.Call("owner.ListAccountAddresses", struct{}{}, &out)
	return
}

func (c *ServiceClient) Release(req ReleaseRequest) (out ReleaseReply, err error) {
	err = c.Call("tx.Release", req, &out)
	return
}

func (c *ServiceClient) Allegation(req AllegationRequest) (out AllegationReply, err error) {
	err = c.Call("tx.Allegation", req, &out)
	return
}

func (c *ServiceClient) Vote(req VoteRequest) (out VoteReply, err error) {
	err = c.Call("tx.Vote", req, &out)
	return
}

func (c *ServiceClient) Stake(req StakeRequest) (out StakeReply, err error) {
	err = c.Call("tx.Stake", req, &out)
	return
}

func (c *ServiceClient) Unstake(req UnstakeRequest) (out UnstakeReply, err error) {
	err = c.Call("tx.Unstake", req, &out)
	return
}

func (c *ServiceClient) Withdraw(req WithdrawRequest) (out WithdrawReply, err error) {
	err = c.Call("tx.Withdraw", req, &out)
	return
}

/* ONS */
func (c *ServiceClient) ONS_CreateRawCreate(req ONSCreateRequest) (out CreateTxReply, err error) {
	err = c.Call("tx.ONS_CreateRawCreate", req, &out)
	return
}
func (c *ServiceClient) ONS_CreateRawUpdate(req ONSUpdateRequest) (out CreateTxReply, err error) {
	err = c.Call("tx.ONS_CreateRawUpdate", req, &out)
	return
}
func (c *ServiceClient) ONS_CreateRawSale(req ONSSaleRequest) (out CreateTxReply, err error) {
	err = c.Call("tx.ONS_CreateRawSale", req, &out)
	return
}
func (c *ServiceClient) ONS_CreateRawBuy(req ONSPurchaseRequest) (out CreateTxReply, err error) {
	err = c.Call("tx.ONS_CreateRawBuy", req, &out)
	return
}

func (c *ServiceClient) ONS_CreateRawSend(req ONSSendRequest) (out CreateTxReply, err error) {
	err = c.Call("tx.ONS_CreateRawSend", req, &out)
	return
}

func (c *ServiceClient) CreateRawSend(req SendTxRequest) (out *CreateTxReply, err error) {
	err = c.Call("tx.CreateRawSend", req, &out)
	return
}

func (c *ServiceClient) CreateRawSendPool(req SendPoolTxRequest) (out *CreateTxReply, err error) {
	err = c.Call("tx.CreateRawSendPool", req, &out)
	return
}

/* Governance */
func (c *ServiceClient) VoteProposal(req VoteProposalRequest) (out *VoteProposalReply, err error) {
	err = c.Call("tx.VoteProposal", req, &out)
	return
}

func (c *ServiceClient) VoteRequests(req VoteRequestRequest) (out VoteRequestReply, err error) {
	err = c.Call("query.VoteRequests", req, &out)
	return
}

func (c *ServiceClient) ListValidators() (out ListValidatorsReply, err error) {
	err = c.Call("query.ListValidators", struct{}{}, &out)
	return
}

func (c *ServiceClient) ListWitnesses(req ListWitnessesRequest) (out ListWitnessesReply, err error) {
	err = c.Call("query.ListWitnesses", req, &out)
	return
}

func (c *ServiceClient) ListCurrencies() (out *ListCurrenciesReply, err error) {
	err = c.Call("query.ListCurrencies", struct{}{}, &out)
	return
}

func (c *ServiceClient) ListProposal(req ListProposalRequest) (out *ListProposalsReply, err error) {
	err = c.Call("query.ListProposal", req, &out)
	return
}

func (c *ServiceClient) ListProposals(req ListProposalsRequest) (out *ListProposalsReply, err error) {
	err = c.Call("query.ListProposals", req, &out)
	return
}

func (c *ServiceClient) ListRewards(req RewardsRequest) (out *ListRewardsReply, err error) {
	err = c.Call("query.ListRewardsForValidator", req, &out)
	return
}

func (c *ServiceClient) WithdrawRewards(req WithdrawRewardsRequest) (out WithdrawRewardsReply, err error) {
	err = c.Call("tx.WithdrawRewards", req, &out)
	return
}

/* Broadcast */
func (c *ServiceClient) TxAsync(req BroadcastRequest) (out BroadcastReply, err error) {
	err = c.Call("broadcast.TxAsync", req, &out)
	return

}

func (c *ServiceClient) TxSync(req BroadcastRequest) (out BroadcastReply, err error) {
	err = c.Call("broadcast.TxSync", req, &out)
	return
}

func (c *ServiceClient) TxCommit(req BroadcastRequest) (out BroadcastReply, err error) {
	err = c.Call("broadcast.TxCommit", req, &out)
	return
}

func (c *ServiceClient) GetTracker(name string) (out BTCGetTrackerReply, err error) {
	/*if len(request) <= 20 {
		return out, errors.New("address has insufficient length")
	}*/
	request := BTCGetTrackerRequest{name}
	err = c.Call("btc.GetTracker", &request, &out)
	return
}

func (c *ServiceClient) CheckCommitResult(hash string, prove bool) (reply TxResponse, err error) {
	request := &TxRequest{Hash: hash, Prove: prove}

	err = c.Call("query.Tx", request, &reply)

	return
}

func (c *ServiceClient) GetTotalNetwkDelegation(onlyActive int) (reply GetTotalNetwkDelgReply, err error) {
	request := &GetTotalNetwkDelegation{OnlyActive: onlyActive}

	err = c.Call("query.GetTotalNetwkDelegation", request, &reply)

	return
}

func (c *ServiceClient) ListDelegation(delegationAddresses []keys.Address) (reply ListDelegationReply, err error) {
	request := &ListDelegationRequest{DelegationAddresses: delegationAddresses}

	err = c.Call("query.ListDelegation", request, &reply)

	return
}
