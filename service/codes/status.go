/*

 */

package codes

import "fmt"

type ServiceError struct {
	Code int
	Msg  string
}

func (se ServiceError) Error() string {
	return fmt.Sprintf("%d: %s", se.Code, se.Msg)
}

var (
	ErrSerialization = ServiceError{InternalError, "error in serialization"}

	ErrLoadingNodeKey = ServiceError{IOError, "error reading node key file"}
	ErrParsingAddress = ServiceError{ParseError, "error parsing address"}
	ErrChainType      = ServiceError{ConfigurationError, "error getting chain type"}

	ErrAddingAccount   = ServiceError{WalletError, "error adding account to wallet"}
	ErrGettingAccount  = ServiceError{WalletError, "error getting account from wallet"}
	ErrDeletingAccount = ServiceError{WalletError, "error in deleting account"}

	ErrGeneratingAccount = ServiceError{AccountsError, "error generating new account"}
	ErrAccountNotFound   = ServiceError{ResourceNotFoundError, "account doesn't in wallet"}
	ErrSigningError      = ServiceError{InternalError, "error while signing"}
	ErrKeyGeneration     = ServiceError{InternalError, "error generating key pair"}

	// Query errors
	ErrBadAddress     = ServiceError{InvalidParams, "address incorrect"}
	ErrGettingBalance = ServiceError{InternalError, "error  getting balance"}
	ErrListValidators = ServiceError{InternalError, "error getting list of validators"}

	// ONS errors
	ErrBadName        = ServiceError{InvalidParams, "domain name not provided"}
	ErrBadOwner       = ServiceError{InvalidParams, "owner address not provided"}
	ErrDomainNotFound = ServiceError{ResourceNotFoundError, "domain not found"}
	ErrFlagNotSet     = ServiceError{InvalidParams, "onsale flag not set"}

	// Tx errors

)
