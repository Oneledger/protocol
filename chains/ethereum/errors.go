package ethereum

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// TODO: Remove this
	ErrNotImplemented Error = "NOT_IMPLEMENTED"
	ErrNilConfig      Error = "was given nil ethereum configuration"
)
