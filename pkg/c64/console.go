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
	cia1 := &CIA1{console: c}
	cia2 := &CIA2{console: c}
	vic := &VICII{console: c}
	vic.rasterPos = 1
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
