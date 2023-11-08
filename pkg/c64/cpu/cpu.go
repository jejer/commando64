package cpu

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/jejer/commando64/pkg/c64"
	"github.com/jejer/commando64/pkg/c64/clock"
)

const (
	// https://www.c64-wiki.com/wiki/Processor_Status_Register
	// N	V	–	B	D	I	Z	C
	// Each flag has a letter symbol for easier reference. Here are the flags and what they mean:
	// Flag				Sym	Bit	Description
	// Negative			N	b7	Set when an operation results in a negative number
	// Overflow			V	b6	Set when a signed addition or subtraction results in an overflow
	// Unused			—	b5	This bit of the processor status register is not used
	// Break			B	b4	Set when a BRK instruction is executed
	// Decimal Mode		D	b3	When set, certain instructions operate in decimal rather than binary mode
	// Interrupt Mask	I	b2	When set, interrupt requests are ignored
	// Zero				Z	b1	Set when an operation results in a zero
	// Carry			C	b0	Set when an unsigned addition or subtraction results in an overflow
	FlagC uint8 = 1 << iota
	FlagZ
	FlagI
	FlagD
	FlagB
	FlagConstant // unused, always set 1
	FlagV
	FlagN

	StackHigh uint16 = 0x01ff // stack works top down
	StackLow  uint16 = 0x0100

	ResetVector uint16 = 0xfffc
	IRQVector   uint16 = 0xfffe
	NMIVector   uint16 = 0xfffa
)

type CPU struct {
	logger     slog.Logger
	clock      *clock.Clock
	mem        c64.MemoryBus
	pc         uint16
	a, x, y, p uint8 // registers
	sp         uint8 // stack pointer
	cycles     int
	irqCh      <-chan bool
}

func NewCPU(logger slog.Logger, clock *clock.Clock, m c64.MemoryBus, irq <-chan bool) *CPU {
	// https://www.c64-wiki.com/index.php/Reset_(Process)
	return &CPU{
		mem:    m,
		pc:     m.ReadWord(ResetVector),
		irqCh:  irq,
		clock:  clock,
		cycles: 0x6,
		logger: *logger.With("Component", "CPU"),
	}
}

func (cpu *CPU) Reset() {
	cpu.a, cpu.x, cpu.y, cpu.p, cpu.sp = 0, 0, 0, 0, 0
	cpu.pc = cpu.mem.ReadWord(ResetVector)
	cpu.cycles = 0x6
}

func (cpu *CPU) Run() {
	for {
		select {
		case isNMI := <-cpu.irqCh:
			if isNMI {
				cpu.NMI()
			} else {
				cpu.IRQ()
			}
		case q := <-cpu.clock.CPU:
			cpu.cycles -= q
			for cpu.cycles <= 0 {
				cpu.step()
			}
		}
	}
}

var CPU_DEBUG_PRINT = 0

func (cpu *CPU) step() {
	if cpu.pc == 0xe5d4 {
		CPU_DEBUG_PRINT = 1
	}
	instraCode := cpu.fetchOP()
	instruction, exist := Instructions[instraCode]
	if CPU_DEBUG_PRINT == 1 {
		CPU_DEBUG_PRINT = 2
		cpu.logger.Debug(`PC-1|   OP    |A |X |Y |P NV.BDIZC|SP|SD| `)
	}
	if CPU_DEBUG_PRINT == 2 {
		cpu.logger.Debug(fmt.Sprintf("%04x|%s%02x%02x%02x|%02x|%02x|%02x|%02x%08b|%02x|%02x| ", cpu.pc-1, runtime.FuncForPC(reflect.ValueOf(instruction.fn).Pointer()).Name()[36:39], cpu.mem.Read(cpu.pc-1), cpu.mem.Read(cpu.pc), cpu.mem.Read(cpu.pc+1), cpu.a, cpu.x, cpu.y, cpu.p, cpu.p, cpu.sp, cpu.mem.Read(StackLow+uint16(cpu.sp)+1)))
	}
	if !exist {
		cpu.logger.Error("Instruction Unsupported", "instruction", instruction)
		panic(1)
	}
	instruction.fn(cpu, instruction.mode)
	cpu.cycles += int(instruction.cycles)
}

func (cpu *CPU) IRQ() {
	if cpu.hasFlag(FlagI) {
		return
	}
	cpu.interrupt(false, IRQVector)
	cpu.cycles += 7
}

func (cpu *CPU) NMI() {
	cpu.interrupt(false, NMIVector)
	cpu.cycles += 7
}

func (cpu *CPU) interrupt(brk bool, vector uint16) {
	cpu.push(uint8((cpu.pc >> 8) & 0xff))
	cpu.push(uint8(cpu.pc & 0xff))
	if brk {
		cpu.push(cpu.p | FlagB)
	} else {
		cpu.push(cpu.p &^ FlagB)
	}
	cpu.setFlag(FlagI, true)
	cpu.pc = cpu.mem.ReadWord(vector)
}

func (cpu *CPU) fetchOP() byte {
	v := cpu.mem.Read(cpu.pc)
	cpu.pc++
	return v
}

func (cpu *CPU) fetchWord() uint16 {
	v := cpu.mem.ReadWord(cpu.pc)
	cpu.pc += 2
	return v
}

// sp is in uint8 range, overflow/underflow is not possiable
func (cpu *CPU) push(v byte) {
	// cpu.logger.Debug("push", "current sp", cpu.sp, "value", v)
	cpu.mem.Write(StackLow+uint16(cpu.sp), v)
	cpu.sp--
}

func (cpu *CPU) pop() byte {
	// cpu.logger.Debug("pop", "current sp", cpu.sp)
	cpu.sp++
	v := cpu.mem.Read(StackLow + uint16(cpu.sp))
	return v
}

func (cpu *CPU) setFlag(flag uint8, v bool) {
	if v {
		cpu.p |= flag
	} else {
		cpu.p &^= flag
	}
	cpu.p |= FlagConstant
}

func (cpu *CPU) hasFlag(flag uint8) bool {
	return (cpu.p & flag) != 0
}

func (cpu *CPU) loadByte(mode AddressingMode) (byte, uint16) {
	// cpu.logger.Debug("loadByte", "mode", mode)

	var v byte = 0
	var addr uint16 = 0
	switch mode {
	case Accumulator:
		v = cpu.a
	case Immidiate:
		v = cpu.fetchOP()
	case Absolute:
		addr = cpu.fetchWord()
		v = cpu.mem.Read(addr)
	case IndexedAbsoluteX:
		addr = cpu.fetchWord()
		addr += uint16(cpu.x)
		v = cpu.mem.Read(addr)
	case IndexedAbsoluteY:
		addr = cpu.fetchWord()
		addr += uint16(cpu.y)
		v = cpu.mem.Read(addr)
	case Zeropage:
		addr = uint16(cpu.fetchOP())
		v = cpu.mem.Read(addr)
	case IndexedZeropageX:
		addr = uint16(cpu.fetchOP())
		addr += uint16(cpu.x)
		addr = addr & 0xff
		v = cpu.mem.Read(addr)
	case IndexedZeropageY:
		addr = uint16(cpu.fetchOP())
		addr += uint16(cpu.y)
		addr = addr & 0xff
		v = cpu.mem.Read(addr)
	case IndexedIndirectX:
		addr = uint16(cpu.fetchOP())
		addr += uint16(cpu.x)
		addr = cpu.mem.ReadWord(addr & 0xff)
		v = cpu.mem.Read(addr)
	case IndirectIndexedY:
		addr = uint16(cpu.fetchOP())
		addr = cpu.mem.ReadWord(addr)
		addr += uint16(cpu.y)
		v = cpu.mem.Read(addr)
	case AbsoluteIndirect: // only get address
		addr = cpu.fetchWord()
		addr = cpu.mem.ReadWord(addr)
	case Relative:
		addr = cpu.pc + uint16(int8(cpu.fetchOP())) // offset could be nagative
	default:
		cpu.logger.Error("unsupported addressing mode", "mode", mode)
	}

	return v, addr
}
