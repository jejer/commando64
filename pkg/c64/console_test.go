package c64

import (
	"log/slog"
	"os"
	"testing"
)

func TestConsole(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	// logger := slog.Default()

	c := NewConsole(*logger)
	// setup c64 default roms
	c.Memory.Write(CpuPortRegister, LORAM|HIRAM|CHAREN)
	c.Memory.LoadRom("../../test/roms/basic.901226-01.bin", BasicRomAddr, false)
	c.Memory.LoadRom("../../test/roms/kernal.901227-03.bin", KernalRomAddr, false)
	c.Memory.LoadRom("../../test/roms/characters.901225-01.bin", CharsRomAddr, false)
	c.CPU.Reset()
	c.Run()
}
