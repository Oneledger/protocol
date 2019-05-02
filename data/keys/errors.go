package keys

import "errors"

var (
	ErrWrongPublicKeyAdapter  = errors.New("error in asserting to PublicKey Addapter")
	ErrWrongPrivateKeyAdapter = errors.New("error in asserting to PrivateKey Addapter")
)
