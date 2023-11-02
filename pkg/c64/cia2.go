package c64

import (
	"fmt"
	"log/slog"
)

type CIA2 struct {
	console *Console
	logger  slog.Logger

	// https://www.c64-wiki.com/wiki/CIA
	// $DD00 Data Port A, keyboard matrix columns
	dataPortA uint8
	// $DD01 Data Port B, keyboard matrix rows
	dataPortB uint8
	// $DD02 Data Direction Port A, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortADir uint8
	// $DD03 Data Direction Port B, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortBDir uint8
	// $DD04 $DD05 TimerA
	timerA uint16
	// $DD06 $DD07 TimerB
	timerB uint16
	// $DD08 ~ $DD0B Real Time Clock, 0.1s, 1s, 1m, 1h
	rtc [4]uint8
	// $DD0C Serial shift register
	sdr uint8
	// $DD0D Interrupt Control and status
	irqControl       uint8
	irqStatus        uint8
	timerAIRQEnabled bool
	timerAEnabled    bool
	timerACounter    uint16
	timerBIRQEnabled bool
	timerBEnabled    bool
	timerBCounter    uint16
	prevCPUCycles    uint64
	// $DD0E Control Timer A
	timerAControl uint8
	// $DD0F Control Timer B
	timerBControl uint8
}

func NewCIA2(c *Console, logger slog.Logger) *CIA2 {
	cia2 := &CIA2{console: c}
	cia2.logger = *logger.With("Component", "CIA2")
	return cia2
}

func (cia2 *CIA2) Write(addr uint16, v uint8) {
	switch addr {
	case 0xdd00:
		cia2.logger.Info(fmt.Sprintf("Write 0xDD00: %08b", v))
		cia2.dataPortA = v
	case 0xdd01:
		cia2.dataPortB = v
	case 0xdd02:
		cia2.dataPortADir = v
	case 0xdd03:
		cia2.dataPortBDir = v
	case 0xdd04:
		cia2.timerA &= 0xff00
		cia2.timerA |= uint16(v)
	case 0xdd05:
		cia2.timerA &= 0x00ff
		cia2.timerA |= (uint16(v) << 8)
	case 0xdd06:
		cia2.timerB &= 0xff00
		cia2.timerB |= uint16(v)
	case 0xdd07:
		cia2.timerB &= 0x00ff
		cia2.timerB |= (uint16(v) << 8)
	case 0xdd08, 0xdd09, 0xdd0a, 0xdd0b: // TODO: TOD registers
	case 0xdd0c: // serial shift register
	case 0xdd0d:
		cia2.irqControl = v
		if v&0x81 == 0x81 {
			cia2.logger.Debug("TimerA IRQ Enabled")
			cia2.timerAIRQEnabled = true
		}
		if v&0x81 == 0x01 {
			cia2.logger.Debug("TimerA IRQ Disabled")
			cia2.timerAIRQEnabled = false
		}
		if v&0x82 == 0x82 {
			cia2.logger.Debug("TimerB IRQ Enabled")
			cia2.timerBIRQEnabled = true
		}
		if v&0x82 == 0x02 {
			cia2.logger.Debug("TimerB IRQ Disabled")
			cia2.timerBIRQEnabled = false
		}
	case 0xdd0e:
		cia2.timerAControl = v
		if v&0x01 == 1 {
			cia2.logger.Debug("TimerA Enabled")
			cia2.timerAEnabled = true
		} else {
			cia2.logger.Debug("TimerA Disabled")
			cia2.timerAEnabled = false
		}
		if v&0x10 != 0 {
			cia2.logger.Debug("Load TimerA")
			cia2.timerACounter = cia2.timerA
		}
	case 0xdd0f:
		cia2.timerBControl = v
		if v&0x01 == 1 {
			cia2.logger.Debug("TimerB Enabled")
			cia2.timerBEnabled = true
		} else {
			cia2.logger.Debug("TimerB Disabled")
			cia2.timerBEnabled = false
		}
		if v&0x10 != 0 {
			cia2.logger.Debug("Load TimerB")
			cia2.timerBCounter = cia2.timerB
		}
	}
}
func (cia2 *CIA2) Read(addr uint16) uint8 {
	switch addr {
	case 0xdd00:
		return cia2.dataPortA
	case 0xdd01:
		return cia2.dataPortB
	case 0xdd02:
		return cia2.dataPortADir
	case 0xdd03:
		return cia2.dataPortBDir
	case 0xdd04:
		return uint8(cia2.timerACounter & 0x00ff)
	case 0xdd05:
		return uint8((cia2.timerACounter & 0xff00) >> 8)
	case 0xdd06:
		return uint8(cia2.timerBCounter & 0x00ff)
	case 0xdd07:
		return uint8((cia2.timerBCounter & 0xff00) >> 8)
	case 0xdd08, 0xdd09, 0xdd0a, 0xdd0b: // TODO: TOD registers
	case 0xdd0c: // serial shift register
	case 0xdd0d:
		return cia2.irqStatus
	case 0xdd0e:
		return cia2.timerAControl
	case 0xdd0f:
		return cia2.timerBControl
	}
	return 0
}

func (cia2 *CIA2) Step() {
	if cia2.timerAEnabled {
		eclipse := cia2.console.CPU.cycles - cia2.prevCPUCycles
		if eclipse > uint64(cia2.timerACounter) {
			if cia2.timerAIRQEnabled {
				cia2.irqStatus |= 0x81
				cia2.console.CPU.NMI()
			}
		}
		cia2.timerACounter -= uint16(eclipse)
	}
	if cia2.timerBEnabled {
		eclipse := cia2.console.CPU.cycles - cia2.prevCPUCycles
		if eclipse > uint64(cia2.timerBCounter) {
			if cia2.timerBIRQEnabled {
				cia2.irqStatus |= 0x82
				cia2.console.CPU.NMI()
			}
		}
		cia2.timerBCounter -= uint16(eclipse)
	}
	cia2.prevCPUCycles = cia2.console.CPU.cycles
}
