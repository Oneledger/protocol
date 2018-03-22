package types

import (
	"testing"
	"time"

	//"github.com/stretchr/testify/require"

	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/tendermint/types"
)

func examplePrevote() *Vote {
	return exampleVote(VoteTypePrevote)
}

func examplePrecommit() *Vote {
	return exampleVote(VoteTypePrecommit)
}

func exampleVote(t byte) *Vote {
	var stamp, err = time.Parse(wire.RFC3339Millis, "2017-12-25T03:00:01.234Z")
	if err != nil {
		panic(err)
	}

	return &Vote{
		Vote: types.Vote{
			ValidatorAddress: []byte("addr"),
			ValidatorIndex:   56789,
			Height:           12345,
			Round:            2,
			Timestamp:        stamp,
			Type:             t,
		},
		Channel: 1,
	}
}

func TestVoteString(t *testing.T) {
	tc := []struct {
		name string
		in   string
		out  string
	}{
		{"Precommit", examplePrecommit().String(), `Vote{56789:616464720000 12345/02/2(Precommit) 686173680000 {<nil>} @ 2017-12-25T03:00:01.234Z}`},
		{"Prevote", examplePrevote().String(), `Vote{56789:616464720000 12345/02/1(Prevote) 686173680000 {<nil>} @ 2017-12-25T03:00:01.234Z}`},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(st *testing.T) {
			if tt.in != tt.out {
				t.Errorf("Got unexpected string for Proposal. Expected:\n%v\nGot:\n%v", tt.in, tt.out)
			}
		})
	}
}


