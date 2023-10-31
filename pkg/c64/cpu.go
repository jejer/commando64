package c64

import (
	"log/slog"
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
	_
	FlagV
	FlagN

	StackHigh uint16 = 0x01ff // stack works top down
	StackLow  uint16 = 0x0100

	ResetVector uint16 = 0xfffc
	IRQVector   uint16 = 0xfffe
)

type CPU struct {
	mem        *C64MemoryMap
	pc         uint16
	a, x, y, p uint8 // registers
	sp         uint8 // stack pointer
	cycles     uint64
	logger     slog.Logger
}

func NewCPU(logger slog.Logger, m *C64MemoryMap) *CPU {
	// https://www.c64-wiki.com/index.php/Reset_(Process)
	return &CPU{
		mem:    m,
		pc:     ResetVector,
		cycles: 0x6,
		logger: *logger.With("Component", "CPU"),
	}
}

func (cpu *CPU) Step() {
	cpu.logger.Debug("Step", "PC", cpu.pc)

	instraCode := cpu.fetchOP()

	instruction, exist := Instructions[instraCode]
	if !exist {
		cpu.logger.Error("Instruction Unsupported", "instruction", instruction)
	}
	if err := instruction.fn(cpu, instruction.mode); err != nil {
		cpu.logger.Error("Instruction Failed", "instruction", instruction)
	}
	cpu.cycles += uint64(instruction.cycles)
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
	cpu.logger.Debug("push", "current sp", cpu.sp, "value", v)
	cpu.mem.Write(StackHigh-uint16(cpu.sp), v)
	cpu.sp++
}

func (cpu *CPU) pop() byte {
	cpu.logger.Debug("pop", "current sp", cpu.sp)
	cpu.sp--
	v := cpu.mem.Read(StackHigh - uint16(cpu.sp))
	return v
}

func (cpu *CPU) setFlag(flag uint8, v bool) {
	if v {
		cpu.p |= flag
	} else {
		cpu.p &^= flag
	}
}

func (cpu *CPU) hasFlag(flag uint8) bool {
	return (cpu.p & flag) != 0
}

func (cpu *CPU) loadByte(mode AddressingMode) (byte, uint16) {
	cpu.logger.Debug("loadByte", "mode", mode)

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
		v = cpu.mem.Read(addr)
	case IndexedZeropageY:
		addr = uint16(cpu.fetchOP())
		addr += uint16(cpu.y)
		v = cpu.mem.Read(addr)
	case IndexedIndirectX:
		addr = uint16(cpu.fetchOP())
		addr += uint16(cpu.x)
		addr = cpu.mem.ReadWord(addr)
		v = cpu.mem.Read(addr)
	case IndirectIndexedY:
		addr = uint16(cpu.fetchOP())
		addr = uint16(cpu.mem.Read(addr))
		addr += uint16(cpu.y)
		addr = cpu.mem.ReadWord(addr)
		v = cpu.mem.Read(addr)
	case AbsoluteIndirect: // only get address
		addr = cpu.fetchWord()
		addr = cpu.mem.ReadWord(addr)
	case Relative:
		addr = cpu.pc + uint16(cpu.fetchOP())
	default:
		cpu.logger.Error("unsupported addressing mode", "mode", mode)
	}

	return v, addr
}
