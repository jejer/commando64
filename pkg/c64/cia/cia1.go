package cia

import (
	"log/slog"
	"time"

	"github.com/jejer/commando64/pkg/c64"
	"github.com/jejer/commando64/pkg/c64/clock"
)

type CIA1 struct {
	logger       slog.Logger
	peripheralIO c64.PeripheralIO
	irqCh        chan<- bool
	clock        *clock.Clock

	// https://www.c64-wiki.com/wiki/CIA
	// $DC00 Data Port A, keyboard matrix columns
	dataPortA uint8
	// $DC01 Data Port B, keyboard matrix rows
	dataPortB uint8
	// $DC02 Data Direction Port A, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortADir uint8
	// $DC03 Data Direction Port B, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortBDir uint8
	// $DC04 $DC05 TimerA
	timerA uint16
	// $DC06 $DC07 TimerB
	timerB uint16
	// $DC08 ~ $DC0B Real Time Clock, 0.1s, 1s, 1m, 1h
	rtc [4]uint8
	// $DC0C Serial shift register
	sdr uint8
	// $DC0D Interrupt Control and status
	irqControl       uint8
	irqStatus        uint8
	timerAIRQEnabled bool
	timerAEnabled    bool
	timerACounter    uint16
	timerBIRQEnabled bool
	timerBEnabled    bool
	timerBCounter    uint16
	// $DC0E Control Timer A
	timerAControl uint8
	// $DC0F Control Timer B
	timerBControl uint8
}

func NewCIA1(logger slog.Logger, clock *clock.Clock, ch chan<- bool, io c64.PeripheralIO) *CIA1 {
	cia1 := &CIA1{peripheralIO: io, irqCh: ch, clock: clock}
	cia1.logger = *logger.With("Component", "CIA1")
	return cia1
}

func (cia1 *CIA1) Write(addr uint16, v uint8) {
	switch addr {
	case 0xdc00:
		cia1.dataPortA = v
	case 0xdc01:
		cia1.dataPortB = v
	case 0xdc02:
		cia1.dataPortADir = v
	case 0xdc03:
		cia1.dataPortBDir = v
	case 0xdc04:
		cia1.timerA &= 0xff00
		cia1.timerA |= uint16(v)
	case 0xdc05:
		cia1.timerA &= 0x00ff
		cia1.timerA |= (uint16(v) << 8)
	case 0xdc06:
		cia1.timerB &= 0xff00
		cia1.timerB |= uint16(v)
	case 0xdc07:
		cia1.timerB &= 0x00ff
		cia1.timerB |= (uint16(v) << 8)
	case 0xdc08, 0xdc09, 0xdc0a, 0xdc0b: // TODO: TOD registers
	case 0xdc0c: // serial shift register
		cia1.sdr = v
	case 0xdc0d:
		cia1.irqControl = v
		if v&0x81 == 0x81 {
			cia1.logger.Debug("TimerA IRQ Enabled")
			cia1.timerAIRQEnabled = true
		}
		if v&0x81 == 0x01 {
			cia1.logger.Debug("TimerA IRQ Disabled")
			cia1.timerAIRQEnabled = false
		}
		if v&0x82 == 0x82 {
			cia1.logger.Debug("TimerB IRQ Enabled")
			cia1.timerBIRQEnabled = true
		}
		if v&0x82 == 0x02 {
			cia1.logger.Debug("TimerB IRQ Disabled")
			cia1.timerBIRQEnabled = false
		}
	case 0xdc0e:
		cia1.timerAControl = v
		if v&0x01 == 1 {
			cia1.logger.Debug("TimerA Enabled")
			cia1.timerAEnabled = true
		} else {
			cia1.logger.Debug("TimerA Disabled")
			cia1.timerAEnabled = false
		}
		if v&0x10 != 0 {
			cia1.logger.Debug("Load TimerA")
			cia1.timerACounter = cia1.timerA
		}
	case 0xdc0f:
		cia1.timerBControl = v
		if v&0x01 == 1 {
			cia1.logger.Debug("TimerB Enabled")
			cia1.timerBEnabled = true
		} else {
			cia1.logger.Debug("TimerB Disabled")
			cia1.timerBEnabled = false
		}
		if v&0x10 != 0 {
			cia1.logger.Debug("Load TimerB")
			cia1.timerBCounter = cia1.timerB
		}
	}
}
func (cia1 *CIA1) Read(addr uint16) uint8 {
	switch addr {
	case 0xdc00:
	case 0xdc01:
		if cia1.dataPortA == 0xff {
			return 0xff
		}
		if cia1.dataPortA != 0 {
			// https://www.c64-wiki.com/wiki/Keyboard#Hardware
			d := ^cia1.dataPortA
			row := uint8(0)
			for i := 0; i < 8; i++ {
				if d&0x01 != 0 {
					break
				}
				row++
				d >>= 1
			}
			return cia1.peripheralIO.ReadKeyboardMatrix(row)
		}
	case 0xdc02:
		return cia1.dataPortADir
	case 0xdc03:
		return cia1.dataPortBDir
	case 0xdc04:
		return uint8(cia1.timerACounter & 0x00ff)
	case 0xdc05:
		return uint8((cia1.timerACounter & 0xff00) >> 8)
	case 0xdc06:
		return uint8(cia1.timerBCounter & 0x00ff)
	case 0xdc07:
		return uint8((cia1.timerBCounter & 0xff00) >> 8)
	case 0xdc08, 0xdc09, 0xdc0a, 0xdc0b: // TODO: TOD registers
	case 0xdc0c:
		return cia1.sdr
	case 0xdc0d:
		return cia1.irqStatus
	case 0xdc0e:
		return cia1.timerAControl
	case 0xdc0f:
		return cia1.timerBControl
	}
	return 0
}

func (cia1 *CIA1) Run() {
	d := time.Duration(time.Second) / (50 * c64.ScreenLines * 63)
	t := time.NewTicker(d)
	for range t.C {
		cia1.step()
	}
}

func (cia1 *CIA1) step() {
	if cia1.timerAEnabled {
		cia1.timerACounter--
		if cia1.timerACounter == 0 {
			if cia1.timerAIRQEnabled {
				cia1.irqStatus |= 0x81
				go func() { cia1.irqCh <- false }()
			}
			cia1.timerACounter = cia1.timerA
		}
	}
	if cia1.timerBEnabled {
		cia1.timerBCounter--
		if cia1.timerBCounter == 0 {
			if cia1.timerBIRQEnabled {
				cia1.irqStatus |= 0x82
				go func() { cia1.irqCh <- false }()
			}
			cia1.timerBCounter = cia1.timerB
		}
	}
}
