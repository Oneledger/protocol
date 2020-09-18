package farm_data

import "errors"

type (
	BatchID string
	FarmID string
)

func (id BatchID) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("BatchID is empty")
	case len(id) != BATCHIDLENGTH:
		return errors.New("BatchID length is incorrect")
	}
	return nil
}

func (id FarmID) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("FarmID is empty")
	case len(id) != FARMIDLENGTH:
		return errors.New("FarmID length is incorrect")
	}
	return nil
}
