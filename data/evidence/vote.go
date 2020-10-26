package evidence

import (
	"fmt"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/tendermint/tendermint/abci/types"
)

type VoteBlock struct {
	Height    int64
	Addresses []keys.Address
}

type CumulativeVote struct {
	Addresses map[string]int64
}

func (es *EvidenceStore) getVoteBlockKey(height int64) []byte {
	key := []byte(fmt.Sprintf("_svb_%d", height))
	return key
}

func (es *EvidenceStore) GetVoteBlock(height int64) (*VoteBlock, error) {
	dat, err := es.Get(es.getVoteBlockKey(height))
	if err != nil {
		return nil, err
	}
	vb := &VoteBlock{
		Height:    height,
		Addresses: make([]keys.Address, 0),
	}
	if len(dat) == 0 {
		return vb, nil
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, vb)
	if err != nil {
		return nil, err
	}
	return vb, nil
}

func (es *EvidenceStore) SetVoteBlock(height int64, votes []types.VoteInfo) error {
	addresses := make([]keys.Address, 0)
	for i := range votes {
		vote := votes[i]
		valAddress := keys.Address(vote.Validator.Address)
		if vote.GetSignedLastBlock() {
			addresses = append(addresses, valAddress)
		}
	}
	vb := &VoteBlock{
		Height:    height,
		Addresses: addresses,
	}
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(vb)
	if err != nil {
		return err
	}
	return es.Set(es.getVoteBlockKey(height), dat)
}

func (es *EvidenceStore) getCumulativeVote() []byte {
	key := []byte("_scv")
	return key
}

func (es *EvidenceStore) GetCumulativeVote() (*CumulativeVote, error) {
	dat, err := es.Get(es.getCumulativeVote())
	if err != nil {
		return nil, err
	}
	cv := &CumulativeVote{
		Addresses: make(map[string]int64),
	}
	if len(dat) == 0 {
		return cv, nil
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, cv)
	if err != nil {
		return nil, err
	}
	return cv, nil
}

func (es *EvidenceStore) SetCumulativeVote(cv *CumulativeVote, currentHeight int64, yHeight int64) error {
	vb, err := es.GetVoteBlock(currentHeight)
	if err != nil {
		return err
	}
	if currentHeight > yHeight {
		diffHeight := currentHeight - yHeight
		vbi, err := es.GetVoteBlock(diffHeight)
		if err != nil {
			return err
		}
		for _, address := range vbi.Addresses {
			hAddr := address.String()
			if _, ok := cv.Addresses[hAddr]; ok {
				cv.Addresses[hAddr]--
				if res, _ := cv.Addresses[hAddr]; res == 0 {
					delete(cv.Addresses, hAddr)
				}
			}
		}
	}
	for _, address := range vb.Addresses {
		hAddr := address.String()
		cv.Addresses[hAddr]++
	}
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(cv)
	if err != nil {
		return err
	}
	return es.Set(es.getCumulativeVote(), dat)
}
