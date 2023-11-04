package c64

type BasicIO interface {
	Read(addr uint16) byte
	Write(addr uint16, v byte)
}

type MemoryBus interface {
	BasicIO
	ReadWord(addr uint16) uint16
	ReadRom(addr uint16) uint8
	VicRead(addr uint16) uint8
}

type PeripheralIO interface {
	Init()
	EventLoop()
	ReadKeyboardMatrix(row uint8) uint8
	SetFramePixel(x int, y uint16, color uint8)
	RefreshScreen()
}
