package filters

import (
	"log"
	"time"

	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Subscription defines a wrapper for the private subscription
type Subscription struct {
	id        rpc.ID
	typ       filters.Type
	event     string
	created   time.Time
	logsCrit  filters.FilterCriteria
	logs      chan []*ethtypes.Log
	hashes    chan []common.Hash
	headers   chan *rpctypes.Header
	installed chan struct{} // closed when the filter is installed
	eventCh   <-chan coretypes.ResultEvent
	err       chan error
}

// ID returns the underlying subscription RPC identifier.
func (s Subscription) ID() rpc.ID {
	return s.id
}

// Unsubscribe to the current subscription from Tendermint Websocket. It sends an error to the
// subscription error channel if unsubscription fails.
func (s *Subscription) Unsubscribe(es *EventSystem) {
	if err := es.unsubscribe(string(s.ID()), s.event); err != nil {
		s.err <- err
	}

	go func() {
		defer func() {
			log.Println("successfully unsubscribed to event", s.event)
		}()

	uninstallLoop:
		for {
			// write uninstall request and consume logs/hashes. This prevents
			// the eventLoop broadcast method to deadlock when writing to the
			// filter event channel while the subscription loop is waiting for
			// this method to return (and thus not reading these events).
			select {
			case es.uninstall <- s:
				break uninstallLoop
			case <-s.logs:
			case <-s.hashes:
			case <-s.headers:
			}
		}
	}()
}

// Err returns the error channel
func (s *Subscription) Err() <-chan error {
	return s.err
}

// Event returns the tendermint result event channel
func (s *Subscription) Event() <-chan coretypes.ResultEvent {
	return s.eventCh
}
