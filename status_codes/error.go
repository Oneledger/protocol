/*

 */

package status_codes

import "fmt"

/*
	Protocol Error definition
*/
type ProtocolError struct {
	Code int
	Msg  string
}

func (se ProtocolError) Error() string {
	return fmt.Sprintf("%d: %s", se.Code, se.Msg)
}

func (se ProtocolError) ErrorMsg() string {
	return se.Msg
}

func (se ProtocolError) Wrap(err error) *ProtocolError {
	return &ProtocolError{se.Code,
		se.Msg + ": " + err.Error()}
}

func WrapError(err error, code int, msg string) *ProtocolError {
	return &ProtocolError{code, msg + ": " + err.Error()}
}

/*

	Status declarations

*/
var (
	ErrSerialization = ProtocolError{InternalErrorSerialization, "error in serialization"}

	ErrLoadingNodeKey = ProtocolError{IOErrorNodeKey, "error reading node key file"}
	ErrParsingAddress = ProtocolError{ParseErrorAddress, "error parsing address"}
	ErrChainType      = ProtocolError{ConfigurationErrorChainType, "error getting chain type"}

	ErrAddingAccount   = ProtocolError{WalletErrorAddingAccount, "error adding account to wallet"}
	ErrGettingAccount  = ProtocolError{WalletErrorGettingAccount, "error getting account from wallet"}
	ErrDeletingAccount = ProtocolError{WalletErrorDeleteAccount, "error in deleting account"}

	ErrGeneratingAccount = ProtocolError{AccountsErrorGeneratingNewAccount, "error generating new account"}
	ErrAccountNotFound   = ProtocolError{AccountNotFound, "account doesn't in wallet"}
	ErrSigningError      = ProtocolError{InternalErrorSigning, "error while signing"}
	ErrKeyGeneration     = ProtocolError{InternalErrorGeneratingKeyPair, "error generating key pair"}

	// Query errors
	ErrBadAddress     = ProtocolError{IncorrectAddress, "address incorrect"}
	ErrGettingBalance = ProtocolError{InternalErrorGettingBalance, "error  getting balance"}
	ErrListValidators = ProtocolError{InternalErrorListValidators, "error getting list of validators"}

	// ONS errors
	ErrBadName        = ProtocolError{DomainMissing, "domain name not provided"}
	ErrBadOwner       = ProtocolError{OwnerAddressMissing, "owner address not provided"}
	ErrDomainNotFound = ProtocolError{DomainNotFound, "domain not found"}
	ErrFlagNotSet     = ProtocolError{OnSaleFlagNotSet, "onsale flag not set"}

	// Tx errors

)
