package c64

import "log/slog"

const (
	// flags
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
		sp:     0xff,
		cycles: 0x6,
		logger: *logger.With("Component", "CPU"),
	}
}

func (cpu *CPU) Step() {
	cpu.logger.Debug("Step", "PC", cpu.pc)

	instraCode := cpu.mem.Read(cpu.pc)
	cpu.pc++

	instruction := Instructions[instraCode]
	if err := instruction.fn(cpu, instruction.mode); err != nil {
		cpu.logger.Error("Instruction Failed", "instruction", instruction)
	}
	cpu.cycles += uint64(instruction.cycles)
}
