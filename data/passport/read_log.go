package passport

import (
	"fmt"
	"time"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type ReadLog struct {
	Org          TokenTypeID  `json:"org"`
	ReadBy       UserID       `json:"readBy"`
	AdminAddress keys.Address `json:"adminAddress"`
	Person       UserID       `json:"person"`
	Address      keys.Address `json:"address"`
	Test         TestType     `json:"test"`
	ReadAt       string       `json:"readAt"`
}

func NewReadLog(org TokenTypeID, readBy UserID, adminAddress keys.Address,
	person UserID, address keys.Address, test TestType, readAt string) *ReadLog {
	log := &ReadLog{
		Org:          org,
		ReadBy:       readBy,
		AdminAddress: adminAddress,
		Person:       person,
		Address:      address,
		Test:         test,
		ReadAt:       readAt,
	}
	if log.ReadAt == "" {
		log.ReadAt = time.Time{}.Format(time.RFC3339)
	}
	return log
}

func (log *ReadLog) TimeRead() time.Time {
	tm, _ := time.Parse(time.RFC3339, log.ReadAt)
	return tm
}

func (log *ReadLog) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(log)
	if err != nil {
		logger.Error("read log not serializable", err)
		return []byte{}
	}
	return value
}

func (log *ReadLog) FromBytes(msg []byte) (*ReadLog, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, log)
	if err != nil {
		logger.Error("failed to deserialize read log from bytes", err)
		return nil, err
	}
	return log, nil
}

func (log *ReadLog) String() string {
	return fmt.Sprintf("Org= %s, ReadBy= %s, Person= %s, Test=%s, ReadAt=%s",
		log.Org, log.ReadBy, log.Person, log.Test, log.ReadAt)
}
