package serialize

import (
	"errors"
)

var (
	ErrIncorrectChannel = errors.New("incorrect channel name")
	ErrMissingAminoCodec = errors.New("missing amino codec")
	ErrIncorrectAminoCodec = errors.New("incorrect amino codec")
)
