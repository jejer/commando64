package c64

type AddressingMode uint8

const (
	Implied AddressingMode = iota + 1
	IndexedIndirectX
	IndirectIndexedY
	Indirect
	Absolute
	AbsoluteX
	AbsoluteY
	Immidiate
	ZeroPageAddressing
	ZeroPageX
	ZeroPageY
	Accumulator
	Relative
)

type InstraFunc func(cpu *CPU, mode AddressingMode) error

type Instruction struct {
	fn     InstraFunc
	mode   AddressingMode
	cycles uint8
}

var Instructions = map[byte]Instruction{
	0x00: Instruction{BRK, Implied, 7},
}

// BRK Force Break
// Operation:  Forced Interrupt PC + 2 toS P toS
// N Z C I D V
// _ _ _ 1 _ _
func BRK(cpu *CPU, mode AddressingMode) error {
	cpu.logger.Debug("BRK", "mode", mode)
	return nil
}
