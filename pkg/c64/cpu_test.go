package c64

import (
	"log/slog"
	"testing"
)

// var cpu *CPU

// func TestMain(m *testing.M) {
// 	mem := NewC64Memory(*slog.Default())
// 	cpu = NewCPU(*slog.Default(), mem)
// 	os.Exit(m.Run())
// }

// func TestPHP(t *testing.T) {
// 	cpu.p = 0xee
// 	PHP(cpu, Implied)
// 	if 1 == 1 {
// 		t.Errorf("Wrong PC value after RTS operation, got: %d, want: %d.", cpu.pc, 101)
// 	}
// }

func TestCPU(t *testing.T) {
	// opts := &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// }
	// handler := slog.NewTextHandler(os.Stdout, opts)
	// logger := slog.New(handler)
	logger := slog.Default()
	mem := NewC64Memory(*logger)
	cpu := NewCPU(*logger, mem)
	mem.Write(CpuPortRegister, 0x0) // umount c64 roms
	mem.LoadRom("../../test/roms/6502_functional_test.bin", 0x400, true)
	cpu.pc = 0x400
	var pc uint16 = 0
	for {
		if pc == cpu.pc {
			t.Errorf("CPU test failed at 0x%x", pc)
			break
		}
		if cpu.pc == 0x3463 {
			t.Logf("CPU test passed!")
			break
		}
		pc = cpu.pc
		cpu.Step()
	}
}