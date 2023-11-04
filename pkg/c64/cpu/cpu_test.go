package cpu

import (
	"log/slog"
	"testing"

	"github.com/jejer/commando64/pkg/c64/clock"
	"github.com/jejer/commando64/pkg/c64/memory"
)

func TestCPU(t *testing.T) {
	// opts := &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// }
	// handler := slog.NewTextHandler(os.Stdout, opts)
	// logger := slog.New(handler)
	logger := slog.Default()
	mem := memory.NewC64Memory(*logger, nil, nil, nil)
	irqCh := make(chan bool)
	clock := clock.NewClock()
	cpu := NewCPU(*logger, clock, mem, irqCh)
	mem.Write(0x01, 0x0) // umount c64 roms
	mem.LoadRom("../../../test/roms/6502_functional_test.bin", 0x400, true)
	cpu.pc = 0x400
	var pc uint16 = 0
	// go cpu.Run()
	// for i := uint64(0); true; i++ {
	// 	t.Logf("pc: %04x", pc)
	// 	if pc == cpu.pc {
	// 		t.Errorf("CPU test failed at 0x%x, i=%d", pc, i)
	// 		break
	// 	}
	// 	if cpu.pc == 0x3463 {
	// 		t.Logf("CPU test passed!")
	// 		break
	// 	}
	// 	pc = cpu.pc
	// 	clock <- i
	// 	<-time.After(time.Duration(time.Millisecond))
	// }
	for i := uint64(0); ; i++ {
		t.Logf("pc: %04x", pc)
		if pc == cpu.pc {
			t.Errorf("CPU test failed at 0x%x", pc)
			break
		}
		if cpu.pc == 0x3463 {
			t.Logf("CPU test passed! i=%d", i)
			break
		}
		pc = cpu.pc
		cpu.step()
	}
}
