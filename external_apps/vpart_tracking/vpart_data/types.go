package vpart_data

import "github.com/pkg/errors"

type (
	Vin         string
	StockNumber string
)

func (vin Vin) Err() error {
	switch {
	case len(vin) == 0:
		return errors.New("VIN is empty")
	case len(vin) != VINLENGTH:
		return errors.New("VIN length is incorrect, should be 17 characters")
	}
	return nil
}

func (stockNum StockNumber) Err() error {
	switch {
	case len(stockNum) == 0:
		return errors.New("stock number is empty")
	case len(stockNum) != STOCKNUMLENGTH:
		return errors.New("stock number length is incorrect, should be 9 characters")
	}
	return nil
}



