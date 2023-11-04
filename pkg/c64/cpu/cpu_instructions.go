package cpu

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

type InstraFunc func(cpu *CPU, mode AddressingMode)

type Instruction struct {
	fn     InstraFunc
	mode   AddressingMode
	cycles uint8
}

var Instructions = map[byte]Instruction{
	0x00: {BRK, Implied, 7},
	0x01: {ORA, IndexedIndirectX, 6},
	// 0x02: Undefined, omitted below
	0x05: {ORA, Zeropage, 3},
	0x06: {ASL, Zeropage, 5},
	0x08: {PHP, Implied, 3},
	0x09: {ORA, Immidiate, 2},
	0x0a: {ASL, Accumulator, 2},
	0x0d: {ORA, Absolute, 4},
	0x0e: {ASL, Absolute, 6},
	0x10: {BPL, Relative, 2},
	0x11: {ORA, IndirectIndexedY, 5},
	0x15: {ORA, IndexedZeropageX, 4},
	0x16: {ASL, IndexedZeropageX, 6},
	0x18: {CLC, Implied, 2},
	0x19: {ORA, IndexedAbsoluteY, 4},
	0x1d: {ORA, IndexedAbsoluteX, 4},
	0x1e: {ASL, IndexedAbsoluteX, 7},
	0x20: {JSR, Absolute, 6},
	0x21: {AND, IndexedIndirectX, 6},
	0x24: {BIT, Zeropage, 3},
	0x25: {AND, Zeropage, 3},
	0x26: {ROL, Zeropage, 5},
	0x28: {PLP, Implied, 4},
	0x29: {AND, Immidiate, 2},
	0x2a: {ROL, Accumulator, 2},
	0x2c: {BIT, Absolute, 4},
	0x2d: {AND, Absolute, 4},
	0x2e: {ROL, Absolute, 6},
	0x30: {BMI, Relative, 2},
	0x31: {AND, IndirectIndexedY, 5},
	0x35: {AND, IndexedZeropageX, 4},
	0x36: {ROL, IndexedZeropageX, 6},
	0x38: {SEC, Implied, 2},
	0x39: {AND, IndexedAbsoluteY, 4},
	0x3d: {AND, IndexedAbsoluteX, 4},
	0x3e: {ROL, IndexedAbsoluteX, 7},
	0x40: {RTI, Implied, 6},
	0x41: {EOR, IndexedIndirectX, 6},
	0x45: {EOR, Zeropage, 3},
	0x46: {LSR, Zeropage, 5},
	0x48: {PHA, Implied, 3},
	0x49: {EOR, Immidiate, 2},
	0x4a: {LSR, Accumulator, 2},
	0x4c: {JMP, Absolute, 3},
	0x4d: {EOR, Absolute, 4},
	0x4e: {LSR, Absolute, 6},
	0x50: {BVC, Relative, 2},
	0x51: {EOR, IndirectIndexedY, 5},
	0x55: {EOR, IndexedZeropageX, 4},
	0x56: {LSR, IndexedZeropageX, 6},
	0x58: {CLI, Implied, 2},
	0x59: {EOR, IndexedAbsoluteY, 4},
	0x5d: {EOR, IndexedAbsoluteX, 4},
	0x5e: {LSR, IndexedAbsoluteX, 7},
	0x60: {RTS, Implied, 6},
	0x61: {ADC, IndexedIndirectX, 6},
	0x65: {ADC, Zeropage, 3},
	0x66: {ROR, Zeropage, 5},
	0x68: {PLA, Implied, 4},
	0x69: {ADC, Immidiate, 2},
	0x6a: {ROR, Accumulator, 2},
	0x6c: {JMP, AbsoluteIndirect, 5},
	0x6d: {ADC, Absolute, 4},
	0x6e: {ROR, Absolute, 6},
	0x70: {BVS, Relative, 2},
	0x71: {ADC, IndirectIndexedY, 5},
	0x75: {ADC, IndexedZeropageX, 4},
	0x76: {ROR, IndexedZeropageX, 6},
	0x78: {SEI, Implied, 2},
	0x79: {ADC, IndexedAbsoluteY, 4},
	0x7d: {ADC, IndexedAbsoluteX, 4},
	0x7e: {ROR, IndexedAbsoluteX, 7},
	0x81: {STA, IndexedIndirectX, 6},
	0x84: {STY, Zeropage, 3},
	0x85: {STA, Zeropage, 3},
	0x86: {STX, Zeropage, 3},
	0x88: {DEY, Implied, 2},
	0x8a: {TXA, Implied, 2},
	0x8c: {STY, Absolute, 4},
	0x8d: {STA, Absolute, 4},
	0x8e: {STX, Absolute, 4},
	0x90: {BCC, Relative, 2},
	0x91: {STA, IndirectIndexedY, 6},
	0x94: {STY, IndexedZeropageX, 4},
	0x95: {STA, IndexedZeropageX, 4},
	0x96: {STX, IndexedZeropageY, 4},
	0x98: {TYA, Implied, 2},
	0x99: {STA, IndexedAbsoluteY, 5},
	0x9a: {TXS, Implied, 2},
	0x9d: {STA, IndexedAbsoluteX, 5},
	0xa0: {LDY, Immidiate, 2},
	0xa1: {LDA, IndexedIndirectX, 6},
	0xa2: {LDX, Immidiate, 2},
	0xa4: {LDY, Zeropage, 3},
	0xa5: {LDA, Zeropage, 3},
	0xa6: {LDX, Zeropage, 3},
	0xa8: {TAY, Implied, 2},
	0xa9: {LDA, Immidiate, 2},
	0xaa: {TAX, Implied, 2},
	0xac: {LDY, Absolute, 4},
	0xad: {LDA, Absolute, 4},
	0xae: {LDX, Absolute, 4},
	0xb0: {BCS, Relative, 2},
	0xb1: {LDA, IndirectIndexedY, 5},
	0xb4: {LDY, IndexedZeropageX, 4},
	0xb5: {LDA, IndexedZeropageX, 4},
	0xb6: {LDX, IndexedZeropageY, 4},
	0xb8: {CLV, Implied, 2},
	0xb9: {LDA, IndexedAbsoluteY, 4},
	0xba: {TSX, Implied, 2},
	0xbc: {LDY, IndexedAbsoluteX, 4},
	0xbd: {LDA, IndexedAbsoluteX, 4},
	0xbe: {LDX, IndexedAbsoluteY, 4},
	0xc0: {CPY, Immidiate, 2},
	0xc1: {CMP, IndexedIndirectX, 6},
	0xc4: {CPY, Zeropage, 3},
	0xc5: {CMP, Zeropage, 3},
	0xc6: {DEC, Zeropage, 5},
	0xc8: {INY, Implied, 2},
	0xc9: {CMP, Immidiate, 2},
	0xca: {DEX, Implied, 2},
	0xcc: {CPY, Absolute, 4},
	0xcd: {CMP, Absolute, 4},
	0xce: {DEC, Absolute, 6},
	0xd0: {BNE, Relative, 2},
	0xd1: {CMP, IndirectIndexedY, 5},
	0xd5: {CMP, IndexedZeropageX, 4},
	0xd6: {DEC, IndexedZeropageX, 6},
	0xd8: {CLD, Implied, 2},
	0xd9: {CMP, IndexedAbsoluteY, 4},
	0xdd: {CMP, IndexedAbsoluteX, 4},
	0xde: {DEC, IndexedAbsoluteX, 7},
	0xe0: {CPX, Immidiate, 2},
	0xe1: {SBC, IndexedIndirectX, 6},
	0xe4: {CPX, Zeropage, 3},
	0xe5: {SBC, Zeropage, 3},
	0xe6: {INC, Zeropage, 5},
	0xe8: {INX, Implied, 2},
	0xe9: {SBC, Immidiate, 2},
	0xea: {NOP, Implied, 2},
	0xec: {CPX, Absolute, 4},
	0xed: {SBC, Absolute, 4},
	0xee: {INC, Absolute, 6},
	0xf0: {BEQ, Relative, 2},
	0xf1: {SBC, IndirectIndexedY, 5},
	0xf5: {SBC, IndexedZeropageX, 4},
	0xf6: {INC, IndexedZeropageX, 6},
	0xf8: {SED, Implied, 2},
	0xf9: {SBC, IndexedAbsoluteY, 4},
	0xfd: {SBC, IndexedAbsoluteX, 4},
	0xfe: {INC, IndexedAbsoluteX, 7},
}

// BRK Force Break
// Operation:  Forced Interrupt
// N Z C I D V
// _ _ _ 1 _ _
func BRK(cpu *CPU, mode AddressingMode) {
	cpu.pc++
	cpu.interrupt(true, IRQVector)
}

// ORA "OR" memory with accumulator
// Operation: A V M -> A
// N Z C I D V
// * * _ _ _ _
func ORA(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.a = cpu.a | v
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, cpu.a&0x80 != 0)
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
func ASL(cpu *CPU, mode AddressingMode) {
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
}

// PHP Push processor status on stack
// Operation:  P toS
// N Z C I D V
// _ _ _ _ _ _
func PHP(cpu *CPU, mode AddressingMode) {
	// cpu.setFlag(FlagB, true) // PHP push the bcf flag active
	cpu.push(cpu.p | FlagB | FlagConstant)
}

// BPL Branch on result plus
// Operation:  Branch on N = 0
// N Z C I D V
// _ _ _ _ _ _
func BPL(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagN) {
		cpu.pc = addr
	}
}

// CLC Clear carry flag
// Operation:  0 -> C
// N Z C I D V
// _ _ 0 _ _ _
func CLC(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagC, false)
}

// JSR Jump to new location saving return address
// Operation:  PC + 2 toS, (PC + 1) -> PCL
//
//	(PC + 2) -> PCH
//
// N Z C I D V
// _ _ _ _ _ _
func JSR(cpu *CPU, mode AddressingMode) {
	addr := cpu.fetchWord()
	cpu.push(uint8((((cpu.pc - 1) >> 8) & 0xff)))
	cpu.push(uint8(((cpu.pc - 1) & 0xff)))
	cpu.pc = addr
}

// AND "AND" memory with accumulator
// Operation:  A /\ M -> A
// N Z C I D V
// * * _ _ _
func AND(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.a = v & cpu.a
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, cpu.a&0x80 != 0)
}

// BIT Test bits in memory with accumulator
// Operation:  A /\ M, M7 -> N, M6 -> V
// Bit 6 and 7 are transferred to the status register.
// If the result of A /\ M is zero then Z = 1, otherwise Z = 0
// N  Z C I D V
// M7 * _ _ _ M6
func BIT(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagV, v&0x40 != 0)
	cpu.setFlag(FlagZ, v&cpu.a == 0)
	cpu.setFlag(FlagN, v&0x80 != 0)
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
func ROL(cpu *CPU, mode AddressingMode) {
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
}

// PLP Pull processor status from stack
// Operation:  P fromS
// From Stack
// _ _ _ _ _ _
func PLP(cpu *CPU, mode AddressingMode) {
	cpu.p = cpu.pop() | FlagConstant
}

// BMI Branch on result minus
// Operation:  Branch on N = 1
// N Z C I D V
// _ _ _ _ _ _
func BMI(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagN) {
		cpu.pc = addr
	}
}

// SEC Set carry flag
// Operation:  1 -> C
// N Z C I D V
// _ _ 1 _ _ _
func SEC(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagC, true)
}

// RTI Return from interrupt
// Operation:  P fromS PC fromS
// N Z C I D V
// * * * * * *
func RTI(cpu *CPU, mode AddressingMode) {
	cpu.p = cpu.pop()
	pc := uint16(cpu.pop())
	pc = uint16(cpu.pop())<<8 + pc
	cpu.pc = pc
}

// EOR "Exclusive-Or" memory with accumulator
// Operation:  A EOR M -> A
// N Z C I D V
// * * _ _ _ _
func EOR(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.a = v ^ cpu.a
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
}

// LSR Shift right one bit (memory or accumulator)
//
//	                 +-+-+-+-+-+-+-+-+
//	Operation:  0 -> |7|6|5|4|3|2|1|0| -> C
//	                 +-+-+-+-+-+-+-+-+
//	N Z C I D V
//	0 * * _ _ _
func LSR(cpu *CPU, mode AddressingMode) {
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
}

// PHA Push accumulator on stack
// Operation:  A toS
// N Z C I D V
// _ _ _ _ _ _
func PHA(cpu *CPU, mode AddressingMode) {
	cpu.push(cpu.a)
}

// JMP Jump to new location
// Operation:  (PC + 1) -> PCL
//
//	(PC + 2) -> PCH
//
// N Z C I D V
// _ _ _ _ _ _
func JMP(cpu *CPU, mode AddressingMode) {
	switch mode {
	case Absolute, AbsoluteIndirect:
		_, addr := cpu.loadByte(mode)
		cpu.pc = addr
	default:
		cpu.logger.Error("unsupported addressing mode for JMP", "mode", mode)
	}
}

// BVC Branch on overflow clear
// Operation:  Branch on V = 0
// N Z C I D V
// _ _ _ _ _ _
func BVC(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagV) {
		cpu.pc = addr
	}
}

// CLI Clear interrupt disable bit
// Operation: 0 -> I
// N Z C I D V
// _ _ _ 0 _ _
func CLI(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagI, false)
}

// RTS Return from Subroutine
// Operation:  PC fromS, PC + 1 -> PC
//
//	N Z C I D V
//	_ _ _ _ _ _
func RTS(cpu *CPU, mode AddressingMode) {
	pc := uint16(cpu.pop())
	pc = uint16(cpu.pop())<<8 + pc
	cpu.pc = pc + 1
}

// ADC Add memory to accumulator with carry
// Operation:  A + M + C -> A, C
// N Z C I D V
// * * * _ _ *
func ADC(cpu *CPU, mode AddressingMode) {
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
func ROR(cpu *CPU, mode AddressingMode) {
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
}

// PLA Pull accumulator from stack
// Operation:  A fromS
// N Z C I D V
// * * _ _ _ _
func PLA(cpu *CPU, mode AddressingMode) {
	cpu.a = cpu.pop()
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
}

// BVS Branch on overflow set
// Operation:  Branch on V = 1
// N Z C I D V
// _ _ _ _ _ _
func BVS(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagV) {
		cpu.pc = addr
	}
}

// SEI Set interrupt disable status
// Operation:  1 -> I
// N Z C I D V
// _ _ _ 1 _ _
func SEI(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagI, true)
}

// STA Store accumulator in memory
// Operation:  A -> M
// N Z C I D V
// _ _ _ _ _ _
func STA(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	cpu.mem.Write(addr, cpu.a)
}

// STY Store Index Y in memory
// Operation:  Y -> M
// N Z C I D V
// _ _ _ _ _ _
func STY(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	cpu.mem.Write(addr, cpu.y)
}

// STX Store Index X in memory
// Operation:  X -> M
// N Z C I D V
// _ _ _ _ _ _
func STX(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	cpu.mem.Write(addr, cpu.x)
}

// DEY Decrement index Y by one
// Operation:  Y - 1 -> Y
// N Z C I D V
// * * _ _ _ _
func DEY(cpu *CPU, mode AddressingMode) {
	cpu.y--
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
}

// TXA Transfer index X to accumulator
// Operation:  X -> A
// N Z C I D V
// * * _ _ _ _
func TXA(cpu *CPU, mode AddressingMode) {
	cpu.a = cpu.x
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
}

// BCC Branch on Carry Clear
// Operation:  Branch on C = 0
// N Z C I D V
// _ _ _ _ _ _
func BCC(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagC) {
		cpu.pc = addr
	}
}

// TYA Transfer index Y to accumulator
// Operation:  Y -> A
// N Z C I D V
// * * _ _ _ _
func TYA(cpu *CPU, mode AddressingMode) {
	cpu.a = cpu.y
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
}

// TXS Transfer index X to stack pointer
// Operation:  X -> S
// N Z C I D V
// _ _ _ _ _ _
func TXS(cpu *CPU, mode AddressingMode) {
	cpu.sp = cpu.x
}

// LDY Load Index Y with memory
// Operation:  M -> Y
// N Z C I D V
// * * _ _ _ _
func LDY(cpu *CPU, mode AddressingMode) {
	cpu.y, _ = cpu.loadByte(mode)
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
}

// LDA Load accumulator with memory
// Operation:  M -> A
// N Z C I D V
// * * _ _ _ _
func LDA(cpu *CPU, mode AddressingMode) {
	cpu.a, _ = cpu.loadByte(mode)
	cpu.setFlag(FlagN, (cpu.a&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.a == 0)
}

// LDX Load Index X with memory
// Operation:  M -> X
// N Z C I D V
// * * _ _ _ _
func LDX(cpu *CPU, mode AddressingMode) {
	cpu.x, _ = cpu.loadByte(mode)
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
}

// TAY Transfer accumulator to index Y
// Operation:  A -> Y
// N Z C I D V
// * * _ _ _ _
func TAY(cpu *CPU, mode AddressingMode) {
	cpu.y = cpu.a
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
}

// TAX Transfer accumulator to index X
// Operation:  A -> X
// N Z C I D V
// * * _ _ _ _
func TAX(cpu *CPU, mode AddressingMode) {
	cpu.x = cpu.a
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
}

// BCS Branch on carry set
// Operation:  Branch on C = 1
// N Z C I D V
// _ _ _ _ _ _
func BCS(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagC) {
		cpu.pc = addr
	}
}

// CLV Clear overflow flag
// Operation: 0 -> V
// N Z C I D V
// _ _ _ _ _ 0
func CLV(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagV, false)
}

// TSX Transfer stack pointer to index X
// Operation:  S -> X
// N Z C I D V
// * * _ _ _ _
func TSX(cpu *CPU, mode AddressingMode) {
	cpu.x = cpu.sp
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
}

// CPY Compare memory and index Y
// Operation:  Y - M
// N Z C I D V
// * * * _ _ _
func CPY(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagC, cpu.y >= v)
	diff := cpu.y - v
	cpu.setFlag(FlagN, (diff&0x80) != 0)
	cpu.setFlag(FlagZ, diff == 0)
}

// CMP Compare memory and accumulator
// Operation:  A - M
// N Z C I D V
// * * * _ _ _
func CMP(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagC, cpu.a >= v)
	diff := cpu.a - v
	cpu.setFlag(FlagN, (diff&0x80) != 0)
	cpu.setFlag(FlagZ, diff == 0)
}

// DEC Decrement Memory by one
// Operation:  M - 1 -> M
// N Z C I D V
// * * _ _ _ _
func DEC(cpu *CPU, mode AddressingMode) {
	v, addr := cpu.loadByte(mode)
	v--
	cpu.mem.Write(addr, v)
	cpu.setFlag(FlagN, (v&0x80) != 0)
	cpu.setFlag(FlagZ, v == 0)
}

// INY Increment Index Y by one
// Operation:  Y + 1 -> Y
// N Z C I D V
// * * _ _ _ _
func INY(cpu *CPU, mode AddressingMode) {
	cpu.y++
	cpu.setFlag(FlagN, (cpu.y&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.y == 0)
}

// DEX Decrement index X by one
// Operation:  X - 1 -> X
// N Z C I D V
// * * _ _ _ _
func DEX(cpu *CPU, mode AddressingMode) {
	cpu.x--
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
}

// BNE Branch on result not zero
// Operation:  Branch on Z = 0
// N Z C I D V
// _ _ _ _ _ _
func BNE(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if !cpu.hasFlag(FlagZ) {
		cpu.pc = addr
	}
}

// CLD Clear decimal mode
// Operation:  0 -> D
// N A C I D V
// _ _ _ _ 0 _
func CLD(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagD, false)
}

// CPX Compare Memory and Index X
// Operation:  X - M
// N Z C I D V
// * * * _ _ _
func CPX(cpu *CPU, mode AddressingMode) {
	v, _ := cpu.loadByte(mode)
	cpu.setFlag(FlagC, cpu.x >= v)
	diff := cpu.x - v
	cpu.setFlag(FlagN, (diff&0x80) != 0)
	cpu.setFlag(FlagZ, diff == 0)
}

// SBC Subtract memory from accumulator with borrow
// Operation:  A - M - C -> A
// N Z C I D V
// * * * _ _ *
// Note:C = Borrow
func SBC(cpu *CPU, mode AddressingMode) {
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
}

// INC Increment memory by one
// Operation:  M + 1 -> M
// N Z C I D V
// * * _ _ _ _
func INC(cpu *CPU, mode AddressingMode) {
	v, addr := cpu.loadByte(mode)
	v++
	cpu.mem.Write(addr, v)
	cpu.setFlag(FlagN, (v&0x80) != 0)
	cpu.setFlag(FlagZ, v == 0)
}

// INX Increment Index X by one
// Operation:  X + 1 -> X
// N Z C I D V
// * * _ _ _ _
func INX(cpu *CPU, mode AddressingMode) {
	cpu.x++
	cpu.setFlag(FlagN, (cpu.x&0x80) != 0)
	cpu.setFlag(FlagZ, cpu.x == 0)
}

// NOP No operation
// Operation:  No Operation (2 cycles)
// N Z C I D V
// _ _ _ _ _ _
func NOP(cpu *CPU, mode AddressingMode) {
}

// BEQ Branch on result zero
// Operation:  Branch on Z = 1
// N Z C I D V
// _ _ _ _ _ _
func BEQ(cpu *CPU, mode AddressingMode) {
	_, addr := cpu.loadByte(mode)
	if cpu.hasFlag(FlagZ) {
		cpu.pc = addr
	}
}

// SED Set decimal mode
// Operation:  1 -> D
// N Z C I D V
// _ _ _ _ 1 _
func SED(cpu *CPU, mode AddressingMode) {
	cpu.setFlag(FlagD, true)
}
