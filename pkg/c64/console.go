package c64

import (
	"log/slog"
)

type Console struct {
	mem    *C64MemoryMap
	cpu    *CPU
	logger slog.Logger
}

func NewConsole(logger slog.Logger) *Console {
	mem := NewC64Memory(logger)
	cpu := NewCPU(logger, mem)
	return &Console{mem, cpu, *logger.With("Component", "Console")}
}
