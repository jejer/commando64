package main

import (
	"fmt"

	"log/slog"

	"github.com/jejer/commando64/pkg/c64"
)

func main() {
	fmt.Println("Hello Commando C64")

	// opts := &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// }
	// handler := slog.NewTextHandler(os.Stdout, opts)
	// logger := slog.New(handler)
	logger := slog.Default()
	c := c64.NewConsole(*logger)
	// setup c64 default roms
	c.Memory.Write(c64.CpuPortRegister, c64.LORAM|c64.HIRAM|c64.CHAREN)
	c.Memory.LoadRom("test/roms/basic.901226-01.bin", c64.BasicRomAddr, false)
	c.Memory.LoadRom("test/roms/kernal.901227-03.bin", c64.KernalRomAddr, false)
	c.Memory.LoadRom("test/roms/characters.901225-01.bin", c64.CharsRomAddr, false)
	c.CPU.Reset()
	go c.Run()
	c.IO.Run()
}
