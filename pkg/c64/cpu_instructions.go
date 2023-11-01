package c64

// https://c64os.com/post/6502instructions

type AddressingMode uint8

const (
	// https://www.c64-wiki.com/wiki/Addressing_mode
	Implied AddressingMode = iota
	Accumulator
	Immidiate
	Absolute
	IndexedAbsoluteX
	IndexedAbsoluteY
	Zeropage
	IndexedZeropageX
	IndexedZeropageY
	Relative
	AbsoluteIndirect
	IndexedIndirectX
	IndirectIndexedY
)

type InstraFunc func(cpu *CPU, mode AddressingMode) error

type Instruction struct {
	fn     InstraFunc
	mode   AddressingMode
	cycles uint8
}

var Instructions = map[byte]Instruction{
	0x00: Instruction{BRK, Implied, 7},
	0x01: Instruction{ORA, IndexedIndirectX, 6},
	// 0x02: Undefined, omitted below
	0x05: Instruction{ORA, Zeropage, 3},
	0x06: Instruction{ASL, Zeropage, 5},
	0x08: Instruction{PHP, Implied, 3},
	0x09: Instruction{ORA, Immidiate, 2},
	0x0a: Instruction{ASL, Accumulator, 2},
	0x0d: Instruction{ORA, Absolute, 4},
	0x0e: Instruction{ASL, Absolute, 6},
	0x10: Instruction{BPL, Relative, 2},
	0x11: Instruction{ORA, IndirectIndexedY, 5},
	0x15: Instruction{ORA, IndexedZeropageX, 4},
	0x16: Instruction{ASL, IndexedZeropageX, 6},
	0x18: Instruction{CLC, Implied, 2},
	0x19: Instruction{ORA, IndexedAbsoluteY, 4},
	0x1d: Instruction{ORA, IndexedAbsoluteX, 4},
	0x1e: Instruction{ASL, IndexedAbsoluteX, 7},
	0x20: Instruction{JSR, Absolute, 6},
	0x21: Instruction{AND, IndexedIndirectX, 6},
	0x24: Instruction{BIT, Zeropage, 3},
	0x25: Instruction{AND, Zeropage, 3},
	0x26: Instruction{ROL, Zeropage, 5},
	0x28: Instruction{PLP, Implied, 4},
	0x29: Instruction{AND, Immidiate, 2},
	0x2a: Instruction{ROL, Accumulator, 2},
	0x2c: Instruction{BIT, Absolute, 4},
	0x2d: Instruction{AND, Absolute, 4},
	0x2e: Instruction{ROL, Absolute, 6},
	0x30: Instruction{BMI, Relative, 2},
	0x31: Instruction{AND, IndirectIndexedY, 5},
	0x35: Instruction{AND, IndexedZeropageX, 4},
	0x36: Instruction{ROL, IndexedZeropageX, 6},
	0x38: Instruction{SEC, Implied, 2},
	0x39: Instruction{AND, IndexedAbsoluteY, 4},
	0x3d: Instruction{AND, IndexedAbsoluteX, 4},
	0x3e: Instruction{ROL, IndexedAbsoluteX, 7},
	0x40: Instruction{RTI, Implied, 6},
	0x41: Instruction{EOR, IndexedIndirectX, 6},
	0x45: Instruction{EOR, Zeropage, 3},
	0x46: Instruction{LSR, Zeropage, 5},
	0x48: Instruction{PHA, Implied, 3},
	0x49: Instruction{EOR, Immidiate, 2},
	0x4a: Instruction{LSR, Accumulator, 2},
	0x4c: Instruction{JMP, Absolute, 3},
	0x4d: Instruction{EOR, Absolute, 4},
	0x4e: Instruction{LSR, Absolute, 6},
	0x50: Instruction{BVC, Relative, 2},
	0x51: Instruction{EOR, IndirectIndexedY, 5},
	0x55: Instruction{EOR, IndexedZeropageX, 4},
	0x56: Instruction{LSR, IndexedZeropageX, 6},
	0x58: Instruction{CLI, Implied, 2},
	0x59: Instruction{EOR, IndexedAbsoluteY, 4},
	0x5d: Instruction{EOR, IndexedAbsoluteX, 4},
	0x5e: Instruction{LSR, IndexedAbsoluteX, 7},
	0x60: Instruction{RTS, Implied, 6},
	0x61: Instruction{ADC, IndexedIndirectX, 6},
	0x65: Instruction{ADC, Zeropage, 3},
	0x66: Instruction{ROR, Zeropage, 5},
	0x68: Instruction{PLA, Implied, 4},
	0x69: Instruction{ADC, Immidiate, 2},
	0x6a: Instruction{ROR, Accumulator, 2},
	0x6c: Instruction{JMP, AbsoluteIndirect, 5},
	0x6d: Instruction{ADC, Absolute, 4},
	0x6e: Instruction{ROR, Absolute, 6},
	0x70: Instruction{BVS, Relative, 2},
	0x71: Instruction{ADC, IndirectIndexedY, 5},
	0x75: Instruction{ADC, IndexedZeropageX, 4},
	0x76: Instruction{ROR, IndexedZeropageX, 6},
	0x78: Instruction{SEI, Implied, 2},
	0x79: Instruction{ADC, IndexedAbsoluteY, 4},
	0x7d: Instruction{ADC, IndexedAbsoluteX, 4},
	0x7e: Instruction{ROR, IndexedAbsoluteX, 7},
	0x81: Instruction{STA, IndexedIndirectX, 6},
	0x84: Instruction{STY, Zeropage, 3},
	0x85: Instruction{STA, Zeropage, 3},
	0x86: Instruction{STX, Zeropage, 3},
	0x88: Instruction{DEY, Implied, 2},
	0x8a: Instruction{TXA, Implied, 2},
	0x8c: Instruction{STY, Absolute, 4},
	0x8d: Instruction{STA, Absolute, 4},
	0x8e: Instruction{STX, Absolute, 4},
	0x90: Instruction{BCC, Relative, 2},
	0x91: Instruction{STA, IndirectIndexedY, 6},
	0x94: Instruction{STY, IndexedZeropageX, 4},
	0x95: Instruction{STA, IndexedZeropageX, 4},
	0x96: Instruction{STX, IndexedZeropageY, 4},
	0x98: Instruction{TYA, Implied, 2},
	0x99: Instruction{STA, IndexedAbsoluteY, 5},
	0x9a: Instruction{TXS, Implied, 2},
	0x9d: Instruction{STA, IndexedAbsoluteX, 5},
	0xa0: Instruction{LDY, Immidiate, 2},
	0xa1: Instruction{LDA, IndexedIndirectX, 6},
	0xa2: Instruction{LDX, Immidiate, 2},
	0xa4: Instruction{LDY, Zeropage, 3},
	0xa5: Instruction{LDA, Zeropage, 3},
	0xa6: Instruction{LDX, Zeropage, 3},
	0xa8: Instruction{TAY, Implied, 2},
	0xa9: Instruction{LDA, Immidiate, 2},
	0xaa: Instruction{TAX, Implied, 2},
	0xac: Instruction{LDY, Absolute, 4},
	0xad: Instruction{LDA, Absolute, 4},
	0xae: Instruction{LDX, Absolute, 4},
	0xb0: Instruction{BCS, Relative, 2},
	0xb1: Instruction{LDA, IndirectIndexedY, 5},
	0xb4: Instruction{LDY, IndexedZeropageX, 4},
	0xb5: Instruction{LDA, IndexedZeropageX, 4},
	0xb6: Instruction{LDX, IndexedZeropageY, 4},
	0xb8: Instruction{CLV, Implied, 2},
	0xb9: Instruction{LDA, IndexedAbsoluteY, 4},
	0xba: Instruction{TSX, Implied, 2},
	0xbc: Instruction{LDY, IndexedAbsoluteX, 4},
	0xbd: Instruction{LDA, IndexedAbsoluteX, 4},
	0xbe: Instruction{LDX, IndexedAbsoluteY, 4},
	0xc0: Instruction{CPY, Immidiate, 2},
	0xc1: Instruction{CMP, IndexedIndirectX, 6},
	0xc4: Instruction{CPY, Zeropage, 3},
	0xc5: Instruction{CMP, Zeropage, 3},
	0xc6: Instruction{DEC, Zeropage, 5},
	0xc8: Instruction{INY, Implied, 2},
	0xc9: Instruction{CMP, Immidiate, 2},
	0xca: Instruction{DEX, Implied, 2},
	0xcc: Instruction{CPY, Absolute, 4},
	0xcd: Instruction{CMP, Absolute, 4},
	0xce: Instruction{DEC, Absolute, 6},
	0xd0: Instruction{BNE, Relative, 2},
	0xd1: Instruction{CMP, IndirectIndexedY, 5},
	0xd5: Instruction{CMP, IndexedZeropageX, 4},
	0xd6: Instruction{DEC, IndexedZeropageX, 6},
	0xd8: Instruction{CLD, Implied, 2},
	0xd9: Instruction{CMP, IndexedAbsoluteY, 4},
	0xdd: Instruction{CMP, IndexedAbsoluteX, 4},
	0xde: Instruction{DEC, IndexedAbsoluteX, 7},
	0xe0: Instruction{CPX, Immidiate, 2},
	0xe1: Instruction{SBC, IndexedIndirectX, 6},
	0xe4: Instruction{CPX, Zeropage, 3},
	0xe5: Instruction{SBC, Zeropage, 3},
	0xe6: Instruction{INC, Zeropage, 5},
	0xe8: Instruction{INX, Implied, 2},
	0xe9: Instruction{SBC, Immidiate, 2},
	0xea: Instruction{NOP, Implied, 2},
	0xec: Instruction{CPX, Absolute, 4},
	0xed: Instruction{SBC, Absolute, 4},
	0xee: Instruction{INC, Absolute, 6},
	0xf0: Instruction{BEQ, Relative, 2},
	0xf1: Instruction{SBC, IndirectIndexedY, 5},
	0xf5: Instruction{SBC, IndexedZeropageX, 4},
	0xf6: Instruction{INC, IndexedZeropageX, 6},
	0xf8: Instruction{SED, Implied, 2},
	0xf9: Instruction{SBC, IndexedAbsoluteY, 4},
	0xfd: Instruction{SBC, IndexedAbsoluteX, 4},
	0xfe: Instruction{INC, IndexedAbsoluteX, 7},
}

// BRK Force Break
// Operation:  Forced Interrupt
// N Z C I D V
// _ _ _ 1 _ _
func BRK(cpu *CPU, mode AddressingMode) error {
	cpu.pc++
	cpu.interrupt(true, IRQVector)
	return nil
}

// ORA "OR" memory with accumulator
// Operation: A V M -> A
// N Z C I D V
// * * _ _ _ _
func ORA(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.a = cpu.a | v
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, cpu.a&0x80 != 0)
	return nil
}

// ASL Shift Left One Bit (Memory or Accumulator)
//
//	+-+-+-+-+-+-+-+-+
//
// Operation:  C <- |7|6|5|4|3|2|1|0| <- 0
//
//	+-+-+-+-+-+-+-+-+
//
// N Z C I D V
// * * * _ _ _
func ASL(cpu *CPU, mode AddressingMode) error {
	asl := func(v uint8) uint8 {
		cpu.setFlag(FlagC, (v&0x80) != 0)
		v = v << 1
		cpu.setFlag(FlagN, (v&0x80) != 0)
		cpu.setFlag(FlagZ, v == 0)
		return v
	}

	switch mode {
	case Accumulator:
		cpu.a = asl(cpu.a)
	case Zeropage, Absolute, IndexedAbsoluteX, IndexedZeropageX:
		v, addr := cpu.loadByte(mode)
		cpu.mem.Write(addr, asl(v))
	default:
		cpu.logger.Error("unsupported addressing mode for ASL", "mode", mode)
	}
	return nil
}

// PHP Push processor status on stack
// Operation:  P toS
// N Z C I D V
// _ _ _ _ _ _
func PHP(cpu *CPU, mode AddressingMode) error {
	// cpu.setFlag(FlagB, true) // PHP push the bcf flag active
	cpu.push(cpu.p | FlagB | FlagConstant)
	return nil
}

// BPL Branch on result plus
// Operation:  Branch on N = 0
// N Z C I D V
// _ _ _ _ _ _
func BPL(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagN) {
		cpu.pc = addr
	}
	return nil
}

// CLC Clear carry flag
// Operation:  0 -> C
// N Z C I D V
// _ _ 0 _ _ _
func CLC(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagC, false)
	return nil
}

// JSR Jump to new location saving return address
// Operation:  PC + 2 toS, (PC + 1) -> PCL
//
//	(PC + 2) -> PCH
//
// N Z C I D V
// _ _ _ _ _ _
func JSR(cpu *CPU, mode AddressingMode) error {
	addr := cpu.fetchWord()
	cpu.push(uint8((((cpu.pc - 1) >> 8) & 0xff)))
	cpu.push(uint8(((cpu.pc - 1) & 0xff)))
	cpu.pc = addr
	return nil
}

// AND "AND" memory with accumulator
// Operation:  A /\ M -> A
// N Z C I D V
// * * _ _ _
func AND(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.a = v & cpu.a
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, cpu.a&0x80 != 0)
	return nil
}

// BIT Test bits in memory with accumulator
// Operation:  A /\ M, M7 -> N, M6 -> V
// Bit 6 and 7 are transferred to the status register.
// If the result of A /\ M is zero then Z = 1, otherwise Z = 0
// N  Z C I D V
// M7 * _ _ _ M6
func BIT(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagV, v&0x40 != 0)
	cpu.setFlag(FlagZ, v&cpu.a == 0)
	cpu.setFlag(FlagN, v&0x80 != 0)
	return nil
}

// ROL Rotate one bit left (memory or accumulator)
//
//	+------------------------------+
//	|         M or A               |
//	|   +-+-+-+-+-+-+-+-+    +-+   |
//
// Operation:   +-< |7|6|5|4|3|2|1|0| <- |C| <-+
//
//	+-+-+-+-+-+-+-+-+    +-+
//
// N Z C I D V
// * * * _ _ _
func ROL(cpu *CPU, mode AddressingMode) error {
	rol := func(v uint8) uint8 {
		var c uint8 = 0
		if cpu.hasFlag(FlagC) {
			c = 1
		}

		cpu.setFlag(FlagC, (v&0x80) != 0)
		v = (v << 1) | c
		cpu.setFlag(FlagN, (v&0x80) != 0)
		cpu.setFlag(FlagZ, v == 0)
		return v
	}

	switch mode {
	case Accumulator:
		cpu.a = rol(cpu.a)
	case Zeropage, Absolute, IndexedAbsoluteX, IndexedZeropageX:
		v, addr := cpu.loadByte(mode)
		cpu.mem.Write(addr, rol(v))
	default:
		cpu.logger.Error("unsupported addressing mode for ROL", "mode", mode)
	}
	return nil
}

// PLP Pull processor status from stack
// Operation:  P fromS
// From Stack
// _ _ _ _ _ _
func PLP(cpu *CPU, mode AddressingMode) error {
	cpu.p = cpu.pop()
	return nil
}

// BMI Branch on result minus
// Operation:  Branch on N = 1
// N Z C I D V
// _ _ _ _ _ _
func BMI(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagN) {
		cpu.pc = addr
	}
	return nil
}

// SEC Set carry flag
// Operation:  1 -> C
// N Z C I D V
// _ _ 1 _ _ _
func SEC(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagC, true)
	return nil
}

// RTI Return from interrupt
// Operation:  P fromS PC fromS
// N Z C I D V
// * * * * * *
func RTI(cpu *CPU, mode AddressingMode) error {
	cpu.p = cpu.pop()
	pc := uint16(cpu.pop())
	pc = uint16(cpu.pop())<<8 + pc
	cpu.pc = pc
	return nil
}

// EOR "Exclusive-Or" memory with accumulator
// Operation:  A EOR M -> A
// N Z C I D V
// * * _ _ _ _
func EOR(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.a = v ^ cpu.a
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// LSR Shift right one bit (memory or accumulator)
//
//	                 +-+-+-+-+-+-+-+-+
//	Operation:  0 -> |7|6|5|4|3|2|1|0| -> C
//	                 +-+-+-+-+-+-+-+-+
//	N Z C I D V
//	0 * * _ _ _
func LSR(cpu *CPU, mode AddressingMode) error {
	lsr := func(v uint8) uint8 {
		cpu.setFlag(FlagC, (v&0x01) != 0)
		v = v >> 1
		cpu.setFlag(FlagN, (v&0x80) != 0)
		cpu.setFlag(FlagZ, v == 0)
		return v
	}

	switch mode {
	case Accumulator:
		cpu.a = lsr(cpu.a)
	case Zeropage, Absolute, IndexedAbsoluteX, IndexedZeropageX:
		v, addr := cpu.loadByte(mode)
		cpu.mem.Write(addr, lsr(v))
	default:
		cpu.logger.Error("unsupported addressing mode for LSR", "mode", mode)
	}
	return nil
}

// PHA Push accumulator on stack
// Operation:  A toS
// N Z C I D V
// _ _ _ _ _ _
func PHA(cpu *CPU, mode AddressingMode) error {
	cpu.push(cpu.a)
	return nil
}

// JMP Jump to new location
// Operation:  (PC + 1) -> PCL
//
//	(PC + 2) -> PCH
//
// N Z C I D V
// _ _ _ _ _ _
func JMP(cpu *CPU, mode AddressingMode) error {
	switch mode {
	case Absolute, AbsoluteIndirect:
		_, addr := cpu.loadByte(mode)
		cpu.pc = addr
	default:
		cpu.logger.Error("unsupported addressing mode for JMP", "mode", mode)
	}
	return nil
}

// BVC Branch on overflow clear
// Operation:  Branch on V = 0
// N Z C I D V
// _ _ _ _ _ _
func BVC(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagV) {
		cpu.pc = addr
	}
	return nil
}

// CLI Clear interrupt disable bit
// Operation: 0 -> I
// N Z C I D V
// _ _ _ 0 _ _
func CLI(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagI, false)
	return nil
}

// RTS Return from Subroutine
// Operation:  PC fromS, PC + 1 -> PC
//
//	N Z C I D V
//	_ _ _ _ _ _
func RTS(cpu *CPU, mode AddressingMode) error {
	pc := uint16(cpu.pop())
	pc = uint16(cpu.pop())<<8 + pc
	cpu.pc = pc + 1
	return nil
}

// ADC Add memory to accumulator with carry
// Operation:  A + M + C -> A, C
// N Z C I D V
// * * * _ _ *
func ADC(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	acc := uint16(cpu.a)
	add := uint16(v)
	var ans uint16 = 0
	var carry uint16 = 0
	if cpu.hasFlag(FlagC) {
		carry = 1
	}

	if cpu.hasFlag(FlagD) {
		// decimal mode
		lo := (acc & 0x0f) + (add & 0x0f) + carry

		var carrylo uint16
		if lo >= 0x0a {
			carrylo = 0x10
			lo -= 0x0a
		}

		hi := (acc & 0xf0) + (add & 0xf0) + carrylo

		if hi >= 0xa0 {
			cpu.setFlag(FlagC, true)
			hi -= 0xa0
		} else {
			cpu.setFlag(FlagC, false)
		}

		ans = hi | lo

		cpu.setFlag(FlagV, ((acc^ans)&0x80) != 0 && ((acc^add)&0x80) == 0)
	} else {
		ans = acc + add + carry
		cpu.setFlag(FlagC, ans > 0xff)
		cpu.setFlag(FlagV, (((acc & 0x80) == (add & 0x80)) && ((acc & 0x80) != (ans & 0x80))))
	}

	cpu.a = uint8(ans)
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// ROR Rotate one bit right (memory or accumulator)
//
//	             +------------------------------+
//	             |                              |
//	             |   +-+    +-+-+-+-+-+-+-+-+   |
//	Operation:   +-> |C| -> |7|6|5|4|3|2|1|0| >-+
//	                 +-+    +-+-+-+-+-+-+-+-+
//
// N Z C I D V
// * * * _ _ _
func ROR(cpu *CPU, mode AddressingMode) error {
	ror := func(v uint8) uint8 {
		var c uint8 = 0
		if cpu.hasFlag(FlagC) {
			c = 1
		}

		cpu.setFlag(FlagC, (v&0x01) != 0)
		v = v>>1 | c<<7
		cpu.setFlag(FlagN, (v&0x80) != 0)
		cpu.setFlag(FlagZ, v == 0)
		return v
	}

	switch mode {
	case Accumulator:
		cpu.a = ror(cpu.a)
	case Zeropage, Absolute, IndexedAbsoluteX, IndexedZeropageX:
		v, addr := cpu.loadByte(mode)
		cpu.mem.Write(addr, ror(v))
	default:
		cpu.logger.Error("unsupported addressing mode for LSR", "mode", mode)
	}
	return nil
}

// PLA Pull accumulator from stack
// Operation:  A fromS
// N Z C I D V
// * * _ _ _ _
func PLA(cpu *CPU, mode AddressingMode) error {
	cpu.a = cpu.pop()
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// BVS Branch on overflow set
// Operation:  Branch on V = 1
// N Z C I D V
// _ _ _ _ _ _
func BVS(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagV) {
		cpu.pc = addr
	}
	return nil
}

// SEI Set interrupt disable status
// Operation:  1 -> I
// N Z C I D V
// _ _ _ 1 _ _
func SEI(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagI, true)
	return nil
}

// STA Store accumulator in memory
// Operation:  A -> M
// N Z C I D V
// _ _ _ _ _ _
func STA(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	cpu.mem.Write(addr, cpu.a)
	return nil
}

// STY Store Index Y in memory
// Operation:  Y -> M
// N Z C I D V
// _ _ _ _ _ _
func STY(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	cpu.mem.Write(addr, cpu.y)
	return nil
}

// STX Store Index X in memory
// Operation:  X -> M
// N Z C I D V
// _ _ _ _ _ _
func STX(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	cpu.mem.Write(addr, cpu.x)
	return nil
}

// DEY Decrement index Y by one
// Operation:  Y - 1 -> Y
// N Z C I D V
// * * _ _ _ _
func DEY(cpu *CPU, mode AddressingMode) error {
	cpu.y--
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
	return nil
}

// TXA Transfer index X to accumulator
// Operation:  X -> A
// N Z C I D V
// * * _ _ _ _
func TXA(cpu *CPU, mode AddressingMode) error {
	cpu.a = cpu.x
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// BCC Branch on Carry Clear
// Operation:  Branch on C = 0
// N Z C I D V
// _ _ _ _ _ _
func BCC(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagC) {
		cpu.pc = addr
	}
	return nil
}

// TYA Transfer index Y to accumulator
// Operation:  Y -> A
// N Z C I D V
// * * _ _ _ _
func TYA(cpu *CPU, mode AddressingMode) error {
	cpu.a = cpu.y
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// TXS Transfer index X to stack pointer
// Operation:  X -> S
// N Z C I D V
// _ _ _ _ _ _
func TXS(cpu *CPU, mode AddressingMode) error {
	cpu.sp = cpu.x
	return nil
}

// LDY Load Index Y with memory
// Operation:  M -> Y
// N Z C I D V
// * * _ _ _ _
func LDY(cpu *CPU, mode AddressingMode) error {
	cpu.y, _ = cpu.loadByte(mode)
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
	return nil
}

// LDA Load accumulator with memory
// Operation:  M -> A
// N Z C I D V
// * * _ _ _ _
func LDA(cpu *CPU, mode AddressingMode) error {
	cpu.a, _ = cpu.loadByte(mode)
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// LDX Load Index X with memory
// Operation:  M -> X
// N Z C I D V
// * * _ _ _ _
func LDX(cpu *CPU, mode AddressingMode) error {
	cpu.x, _ = cpu.loadByte(mode)
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
	return nil
}

// TAY Transfer accumulator to index Y
// Operation:  A -> Y
// N Z C I D V
// * * _ _ _ _
func TAY(cpu *CPU, mode AddressingMode) error {
	cpu.y = cpu.a
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
	return nil
}

// TAX Transfer accumulator to index X
// Operation:  A -> X
// N Z C I D V
// * * _ _ _ _
func TAX(cpu *CPU, mode AddressingMode) error {
	cpu.x = cpu.a
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
	return nil
}

// BCS Branch on carry set
// Operation:  Branch on C = 1
// N Z C I D V
// _ _ _ _ _ _
func BCS(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagC) {
		cpu.pc = addr
	}
	return nil
}

// CLV Clear overflow flag
// Operation: 0 -> V
// N Z C I D V
// _ _ _ _ _ 0
func CLV(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagV, false)
	return nil
}

// TSX Transfer stack pointer to index X
// Operation:  S -> X
// N Z C I D V
// * * _ _ _ _
func TSX(cpu *CPU, mode AddressingMode) error {
	cpu.x = cpu.sp
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
	return nil
}

// CPY Compare memory and index Y
// Operation:  Y - M
// N Z C I D V
// * * * _ _ _
func CPY(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagC, cpu.y >= v)
	diff := cpu.y - v
	cpu.setFlag(FlagN, (diff&0x80) != 0)
	cpu.setFlag(FlagZ, diff == 0)
	return nil
}

// CMP Compare memory and accumulator
// Operation:  A - M
// N Z C I D V
// * * * _ _ _
func CMP(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagC, cpu.a >= v)
	diff := cpu.a - v
	cpu.setFlag(FlagN, (diff&0x80) != 0)
	cpu.setFlag(FlagZ, diff == 0)
	return nil
}

// DEC Decrement Memory by one
// Operation:  M - 1 -> M
// N Z C I D V
// * * _ _ _ _
func DEC(cpu *CPU, mode AddressingMode) error {
	v, addr := cpu.loadByte(mode)
	v--
	cpu.mem.Write(addr, v)
	cpu.setFlag(FlagN, (v&0x80) != 0)
	cpu.setFlag(FlagZ, v == 0)
	return nil
}

// INY Increment Index Y by one
// Operation:  Y + 1 -> Y
// N Z C I D V
// * * _ _ _ _
func INY(cpu *CPU, mode AddressingMode) error {
	cpu.y++
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
	return nil
}

// DEX Decrement index X by one
// Operation:  X - 1 -> X
// N Z C I D V
// * * _ _ _ _
func DEX(cpu *CPU, mode AddressingMode) error {
	cpu.x--
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
	return nil
}

// BNE Branch on result not zero
// Operation:  Branch on Z = 0
// N Z C I D V
// _ _ _ _ _ _
func BNE(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagZ) {
		cpu.pc = addr
	}
	return nil
}

// CLD Clear decimal mode
// Operation:  0 -> D
// N A C I D V
// _ _ _ _ 0 _
func CLD(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagD, false)
	return nil
}

// CPX Compare Memory and Index X
// Operation:  X - M
// N Z C I D V
// * * * _ _ _
func CPX(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagC, cpu.x >= v)
	diff := cpu.x - v
	cpu.setFlag(FlagN, (diff&0x80) != 0)
	cpu.setFlag(FlagZ, diff == 0)
	return nil
}

// SBC Subtract memory from accumulator with borrow
// Operation:  A - M - C -> A
// N Z C I D V
// * * * _ _ *
// Note:C = Borrow
func SBC(cpu *CPU, mode AddressingMode) error {
	v, _ := cpu.loadByte(mode)
	acc := uint16(cpu.a)
	sub := uint16(v)
	var ans uint16 = 0
	var carry uint16 = 0
	if cpu.hasFlag(FlagC) {
		carry = 1
	}

	if cpu.hasFlag(FlagD) {
		lo := 0x0f + (acc & 0x0f) - (sub & 0x0f) + carry

		var carrylo uint16
		if lo < 0x10 {
			lo -= 0x06
			carrylo = 0
		} else {
			lo -= 0x10
			carrylo = 0x10
		}

		hi := 0xf0 + (acc & 0xf0) - (sub & 0xf0) + carrylo

		if hi < 0x100 {
			cpu.setFlag(FlagC, false)
			hi -= 0x60
		} else {
			cpu.setFlag(FlagC, true)
			hi -= 0x100
		}

		ans = hi | lo

		cpu.setFlag(FlagV, ((acc^ans)&0x80) != 0 && ((acc^sub)&0x80) != 0)
	} else {
		ans = 0xff + acc - sub + carry
		cpu.setFlag(FlagC, ans > 0xff)
		cpu.setFlag(FlagV, (((cpu.a & 0x80) != (v & 0x80)) && ((cpu.a & 0x80) != (uint8(ans) & 0x80))))
	}

	cpu.a = uint8(ans)
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
	return nil
}

// INC Increment memory by one
// Operation:  M + 1 -> M
// N Z C I D V
// * * _ _ _ _
func INC(cpu *CPU, mode AddressingMode) error {
	v, addr := cpu.loadByte(mode)
	v++
	cpu.mem.Write(addr, v)
	cpu.setFlag(FlagN, (v&0x80) != 0)
	cpu.setFlag(FlagZ, v == 0)
	return nil
}

// INX Increment Index X by one
// Operation:  X + 1 -> X
// N Z C I D V
// * * _ _ _ _
func INX(cpu *CPU, mode AddressingMode) error {
	cpu.x++
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
	return nil
}

// NOP No operation
// Operation:  No Operation (2 cycles)
// N Z C I D V
// _ _ _ _ _ _
func NOP(cpu *CPU, mode AddressingMode) error {
	return nil
}

// BEQ Branch on result zero
// Operation:  Branch on Z = 1
// N Z C I D V
// _ _ _ _ _ _
func BEQ(cpu *CPU, mode AddressingMode) error {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagZ) {
		cpu.pc = addr
	}
	return nil
}

// SED Set decimal mode
// Operation:  1 -> D
// N Z C I D V
// _ _ _ _ 1 _
func SED(cpu *CPU, mode AddressingMode) error {
	cpu.setFlag(FlagD, true)
	return nil
}
