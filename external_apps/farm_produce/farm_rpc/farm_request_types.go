package farm_rpc

import "github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"

type GetBatchByIDRequest struct {
	BatchID farm_data.BatchID `json:"batchId"`
}

type GetBatchByIDReply struct {
	ProduceBatch farm_data.Produce `json:"produceBatch"`
	Height       int64             `json:"height"`
}
