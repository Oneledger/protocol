package keys

import "errors"

var (
	ErrWrongPublicKeyAdapter  = errors.New("error in asserting to PublicKey Addapter")
	ErrWrongPrivateKeyAdapter = errors.New("error in asserting to PrivateKey Addapter")
	ErrMissMsg                = errors.New("miss message to sign")
	ErrMissSigners            = errors.New("signers not specified")
	ErrInvalidThreshold       = errors.New("invalid threshold")
	ErrNotExpectedSigner      = errors.New("not expected signer")
	ErrInvalidSignedMsg       = errors.New("invalid signed message")
)
