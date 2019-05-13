package app

type Code uint32

const (
	CodeOK    Code = 0
	CodeNotOK Code = 1
)

func (c Code) uint32() uint32 {
	return uint32(c)
}
