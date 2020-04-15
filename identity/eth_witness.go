package identity

import (
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type Witness struct {
	Address     keys.Address   `json:"address"`
	PubKey      keys.PublicKey `json:"pubKey"`
	ECDSAPubKey keys.PublicKey `json:"ecdsaPubkey"`
	Name        string         `json:"name"`
	Chain       chain.Type     `json:"chain"`
}

func (w *Witness) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(w)
	if err != nil {
		logger.Error("ethereum witness not serializable", err)
		return []byte{}
	}
	return value
}

func (w *Witness) FromBytes(msg []byte) (*Witness, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, w)
	if err != nil {
		logger.Error("failed to deserialize ethereum witness from bytes", err)
		return nil, err
	}
	return w, nil
}
