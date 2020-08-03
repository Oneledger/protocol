package evidence

import (
	"errors"
	"fmt"
	"time"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type MaliciousValidators struct {
	Addresses []string
}

type LastValidatorHistory struct {
	Address       keys.Address
	Status        int8
	FrozenHeight  int64
	FrozenAt      *time.Time
	ReleaseHeight int64
	ReleaseAt     *time.Time
}

// TODO
// add cache mechanizm for large block diff count

func (lvh *LastValidatorHistory) Bytes() ([]byte, error) {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(lvh)
	if err != nil {
		return []byte{}, fmt.Errorf("validator not serializable %s", err)
	}
	return value, nil
}

func (lvh *LastValidatorHistory) FromBytes(msg []byte) (*LastValidatorHistory, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, lvh)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize account from bytes %s", err)
	}
	return lvh, nil
}

func (lvh *LastValidatorHistory) IsFrozen() bool {
	if lvh.ReleaseAt == nil {
		return true
	}
	return !lvh.ReleaseAt.After(*lvh.FrozenAt)
}

func (lvh *LastValidatorHistory) ReleaseReady(options *Options, blockCreatedAt time.Time) (bool, error) {
	switch lvh.Status {
	case MISSED_REQUIRED_VOTES:
		return true, nil
	case BYZANTINE_FAULT:
		releaseAt := lvh.FrozenAt.AddDate(0, 0, int(options.ValidatorReleaseTime))
		if !blockCreatedAt.After(releaseAt) {
			return false, fmt.Errorf("Validator could be released after: %s", releaseAt)
		}
		return true, nil
	default:
		return false, errors.New("Unsupported status release")
	}
}

func NewLastValidatorHistory(validatorAddress keys.Address, status int8, height int64, createdAt *time.Time) *LastValidatorHistory {
	return &LastValidatorHistory{
		Address:      validatorAddress,
		Status:       status,
		FrozenHeight: height,
		FrozenAt:     createdAt,
	}
}
