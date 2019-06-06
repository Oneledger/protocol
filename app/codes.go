package app

type Code uint32

const (
	CodeOK    Code = 0x00
	CodeNotOK Code = 0x01
)

func (c Code) uint32() uint32 {
	return uint32(c)
}
