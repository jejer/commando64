package c64

import (
	"log/slog"
)

type Console struct {
	Memory *C64MemoryMap
	CPU    *CPU
	logger slog.Logger
}

func NewConsole(logger slog.Logger) *Console {
	mem := NewC64Memory(logger)
	cpu := NewCPU(logger, mem)
	return &Console{mem, cpu, *logger.With("Component", "Console")}
}

// func (c *Console) LoadRom(path string, addr uint16, ram bool) {
// 	// c.mem.LoadRom()
// }

func (c *Console) Run() {
	for {
		c.CPU.Step()
	}
}
