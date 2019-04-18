package id

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/serialize"
)

type SequenceRecord struct {
	Sequence int64
}

func init() {
	serial.Register(SequenceRecord{})

	serialize.RegisterConcrete(new(SequenceRecord), TagSequenceRecord)
}

func NextSequence(app interface{}, accountkey AccountKey) SequenceRecord {
	sequence := int64(1)
	sequenceDb := GetSequence(app)
	raw := sequenceDb.Get(accountkey)
	if raw != nil {
		interim := raw.(*SequenceRecord)
		sequence = interim.Sequence + 1
	}

	sequenceRecord := SequenceRecord{
		Sequence: sequence,
	}

	session := sequenceDb.Begin()
	session.Set(accountkey, sequenceRecord)
	session.Commit()

	return sequenceRecord
}
