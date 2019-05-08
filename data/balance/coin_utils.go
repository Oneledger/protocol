/*
   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/


	Copyright 2017 - 2019 OneLedger

*/

package balance

import (
	"math/big"
	"strconv"
)

// Handle an incoming string
func parseString(amount string, base *big.Float) *big.Int {
	//logger.Dump("Parsing Amount", amount)
	if amount == "" {
		logger.Error("Empty Amount String", "amount", amount)
		return nil
	}

	/*
		value := new(big.Float)

		_, err := fmt.Sscan(amount, value)
		if err != nil {
			logger.Error("Invalid Float String", "err", err, "amount", amount)
			return nil
		}
		result := bfloat2bint(value, base)
	*/

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Error("Invalid Float String", "err", err, "amount", amount)
		return nil
	}

	result := float2bint(value, base)

	return result
}

func units2bint(amount int64, base *big.Float) *big.Int {
	value := new(big.Int).SetInt64(amount)
	return value
}

// Handle an incoming int (often used for comaparisons)
func int2bint(amount int64, base *big.Float) *big.Int {
	value := new(big.Float).SetInt64(amount)

	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	//logger.Dump("int2bint", amount, result)
	return result
}

// Handle an incoming float
func float2bint(amount float64, base *big.Float) *big.Int {
	value := big.NewFloat(amount)

	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	return result
}
