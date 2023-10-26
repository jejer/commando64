package c64

type Memory interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}
