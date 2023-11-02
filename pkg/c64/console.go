package c64

import (
	"log/slog"
)

type Console struct {
	CIA1   *CIA1
	CIA2   *CIA2
	VIC    *VICII
	Memory *C64MemoryMap
	CPU    *CPU
	logger slog.Logger
}

func NewConsole(logger slog.Logger) *Console {
	c := &Console{logger: *logger.With("Component", "Console")}
	cia1 := NewCIA1(c, logger)
	cia2 := NewCIA2(c, logger)
	vic := NewVICII(c, logger)
	mem := NewC64Memory(c, logger)
	cpu := NewCPU(logger, mem)
	c.CIA1 = cia1
	c.CIA2 = cia2
	c.VIC = vic
	c.Memory = mem
	c.CPU = cpu
	return c
}

// func (c *Console) LoadRom(path string, addr uint16, ram bool) {
// 	// c.mem.LoadRom()
// }

func (c *Console) Run() {
	for {
		c.CPU.Step()
	}
}
