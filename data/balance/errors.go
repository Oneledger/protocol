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
	"errors"

	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrWrongBalanceAdapter = errors.New("error in asserting to BalanceAdapter")
	ErrDuplicateCurrency   = errors.New("provided currency has already been registered")
	ErrMismatchingCurrency = errors.New("mismatching currencies")

	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrBalanceErrorAddFailed   = codes.ProtocolError{Code: codes.BalanceErrorAddFailed, Msg: "Failed to add balance to account"}
	ErrBalanceErrorMinusFailed = codes.ProtocolError{Code: codes.BalanceErrorMinusFailed, Msg: "Failed to minus balance from account"}

	ErrAccountNotFound = errors.New("account not found")
)
