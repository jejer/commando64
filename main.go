package main

import (
	"fmt"

	"log/slog"

	"github.com/jejer/commando64/pkg/c64"
	"github.com/jejer/commando64/pkg/c64/cia"
	"github.com/jejer/commando64/pkg/c64/clock"
	"github.com/jejer/commando64/pkg/c64/cpu"
	"github.com/jejer/commando64/pkg/c64/memory"
	"github.com/jejer/commando64/pkg/c64/peripheral"
	"github.com/jejer/commando64/pkg/c64/vic"
)

func main() {
	fmt.Println("Hello Commando C64")

	// opts := &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// }
	// handler := slog.NewTextHandler(os.Stdout, opts)
	// logger := slog.New(handler)
	logger := slog.Default()

	peripheral := peripheral.NewPeripheralSDL(*logger)
	peripheral.Init()
	irqCh := make(chan bool)
	clock := clock.NewClock()
	cia1 := cia.NewCIA1(*logger, clock, irqCh, peripheral)
	cia2 := cia.NewCIA2(*logger, clock, irqCh)
	memory := memory.NewC64Memory(*logger, cia1, cia2, nil)
	vic := vic.NewVICII(*logger, clock, memory, irqCh, peripheral)
	memory.SetVIC(vic)
	cpu := cpu.NewCPU(*logger, clock, memory, irqCh)

	memory.Write(0x01, 0x07)
	memory.LoadRom("test/roms/basic.901226-01.bin", c64.BasicRomAddr, false)
	memory.LoadRom("test/roms/kernal.901227-03.bin", c64.KernalRomAddr, false)
	memory.LoadRom("test/roms/characters.901225-01.bin", c64.CharsRomAddr, false)

	cpu.Reset()
	go cia1.Run()
	go cia2.Run()
	go cpu.Run()
	go vic.Run()
	go clock.Run()

	peripheral.EventLoop()
}
