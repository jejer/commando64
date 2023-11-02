package c64

import (
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"
)

const (
	// pages
	// https://www.c64-wiki.com/wiki/Memory_Map
	ZeroPage          uint16 = 0x0000
	CartLoStartPage   uint16 = 0x8000
	CartLoEndPage     uint16 = 0x9F00
	BasicStartPage    uint16 = 0xa000
	BasicEndPage      uint16 = 0xbf00
	CartHi1StartPage  uint16 = 0xa000
	CartHi1EndPage    uint16 = 0xbf00
	CharStartPage     uint16 = 0xd000
	CharEndPage       uint16 = 0xdf00
	VICStartPage      uint16 = 0xd000
	VICEndPage        uint16 = 0xd300
	SIDStartPage      uint16 = 0xd400
	SIDEndPage        uint16 = 0xd700
	ColorRamStartPage uint16 = 0xd800
	ColorRamEndPage   uint16 = 0xdb00
	CIA1Page          uint16 = 0xdc00
	CIA2Page          uint16 = 0xdd00
	IO1Page           uint16 = 0xde00
	IO2Page           uint16 = 0xdf00
	KernalStartPage   uint16 = 0xe000
	KernalEndPage     uint16 = 0xff00
	CartHi2StartPage  uint16 = 0xe000
	CartHi2EndPage    uint16 = 0xff00

	// roms
	BasicRomAddr  uint16 = 0xa000
	KernalRomAddr uint16 = 0xe000
	CharsRomAddr  uint16 = 0xd000

	// banking switching
	LORAM  byte = 1 << 0 // BIT0: Configures RAM or ROM at $A000-$BFFF for basic rom
	HIRAM  byte = 1 << 1 // BIT1: Configures RAM or ROM at $E000-$FFFF for kernal rom
	CHAREN byte = 1 << 2 // BIT2: Configures ROM or I/O at $D000-$DFFF for character rom

	// registers
	CpuPortRegister uint16 = 0x0001 // for banking switch
)

type BandMode uint8

const BandModeIO BandMode = 0
const BandModeRAM BandMode = 1
const BandModeROM BandMode = 2

type Memory interface {
	Read(addr uint16) byte
	Write(addr uint16, v byte)
}

type C64MemoryMap struct {
	console *Console
	ram     [65536]byte
	rom     [65536]byte
	logger  slog.Logger
}

func NewC64Memory(c *Console, logger slog.Logger) *C64MemoryMap {
	m := &C64MemoryMap{console: c}
	m.logger = *logger.With("Component", "Memory")
	return m
}

func (m *C64MemoryMap) Write(addr uint16, v byte) {
	// page := addr & 0xff00
	// switch {
	// case page == ZeroPage:
	// 	m.ram[addr] = v
	// 	// log ROM bank switching
	// 	if CpuPortRegister == addr {
	// 		m.RomBankSwitch(v)
	// 	}
	// 	return
	// default:
	// 	// C64 always write to RAM even ROM is mounted.
	// 	m.ram[addr] = v
	// }

	if CpuPortRegister == addr {
		m.ram[addr] = v
		m.RomBankSwitch(v)
		return
	}

	switch m.GetAddrBandMode(addr) {
	case BandModeIO:
		page := addr & 0xff00
		switch {
		case page >= VICStartPage && page <= VICEndPage:
			m.console.VIC.Write(addr, v)
		case page == CIA1Page:
			m.console.CIA1.Write(addr, v)
		case page == CIA2Page:
			m.console.CIA2.Write(addr, v)
		default:
			m.ram[addr] = v
		}
	default:
		// C64 always write to RAM even ROM is mounted.
		m.ram[addr] = v
	}
}

func (m *C64MemoryMap) Read(addr uint16) byte {
	switch m.GetAddrBandMode(addr) {
	case BandModeROM:
		return m.rom[addr]
	case BandModeIO:
		page := addr & 0xff00
		switch {
		case page >= VICStartPage && page <= VICEndPage:
			return m.console.VIC.Read(addr)
		case page == CIA1Page:
			return m.console.CIA1.Read(addr)
		case page == CIA2Page:
			return m.console.CIA2.Read(addr)
		default:
			return m.ram[addr]
		}
	default:
		return m.ram[addr]
	}
	// page := addr & 0xff00
	// switch {
	// case page >= CharStartPage && page <= CharEndPage:
	// 	switch m.GetAddrBandMode(addr) {
	// 	case BandModeROM:
	// 		return m.rom[addr]
	// 	default:
	// 		return m.ram[addr]
	// 	}
	// default:
	// 	return m.ram[addr]
	// }
}

func (m *C64MemoryMap) ReadWord(addr uint16) uint16 {
	return uint16(m.Read(addr)) | (uint16(m.Read(addr+1)) << 8)
}

// https://web.archive.org/web/20230527235630/https://www.c64-wiki.com/wiki/Zeropage
func (m *C64MemoryMap) RomBankSwitch(v byte) {
	m.logger.Info("RomBankSwitch", "data", fmt.Sprintf("%08b", v))
}

func (m *C64MemoryMap) LoadRom(path string, addr uint16, ram bool) error {
	file, err := os.Open(path)
	if err != nil {
		m.logger.Error("LoadRom Can't open file", "path", path, "addr", fmt.Sprintf("%08x", addr), "err", err)
		return err
	}

	byteContent, err := ioutil.ReadAll(file)
	if err != nil {
		m.logger.Error("LoadRom Can't read file", "path", path, "addr", fmt.Sprintf("%08x", addr), "err", err)
		return err
	}

	for i := 0; i < len(byteContent); i++ {
		if ram {
			m.ram[addr+uint16(i)] = byteContent[i]
		} else {
			m.rom[addr+uint16(i)] = byteContent[i]
		}
	}
	return nil
}

func (m *C64MemoryMap) GetAddrBandMode(addr uint16) BandMode {
	// https://web.archive.org/web/20230714080427/https://www.c64-wiki.com/wiki/Bank_Switching
	// https://web.archive.org/web/20201029042742/http://unusedino.de/ec64/technical/aay/c64/memcfg.htm
	//      Bit+-------------+-----------+------------+
	//      210| $8000-$BFFF |$D000-$DFFF|$E000-$FFFF |
	// +---+---+-------------+-----------+------------+
	// | 7 |111| Cart.+Basic |    I/O    | Kernal ROM |
	// +---+---+-------------+-----------+------------+
	// | 6 |110|     RAM     |    I/O    | Kernal ROM |
	// +---+---+-------------+-----------+------------+
	// | 5 |101|     RAM     |    I/O    |    RAM     |
	// +---+---+-------------+-----------+------------+
	// | 4 |100|     RAM     |    RAM    |    RAM     |
	// +---+---+-------------+-----------+------------+
	// | 3 |011| Cart.+Basic | Char. ROM | Kernal ROM |
	// +---+---+-------------+-----------+------------+
	// | 2 |010|     RAM     | Char. ROM | Kernal ROM |
	// +---+---+-------------+-----------+------------+
	// | 1 |001|     RAM     | Char. ROM |    RAM     |
	// +---+---+-------------+-----------+------------+
	// | 0 |000|     RAM     |    RAM    |    RAM     |
	// +---+---+-------------+-----------+------------+
	hiram := ((m.ram[CpuPortRegister] & HIRAM) != 0)   // kernal
	loram := ((m.ram[CpuPortRegister] & LORAM) != 0)   // basic
	charen := ((m.ram[CpuPortRegister] & CHAREN) != 0) // char

	// TODO: support cartridge and expansion cards
	page := addr & 0xff00
	switch {
	case page >= KernalStartPage && page <= KernalEndPage: // 0xe000 ~ 0xffff
		if hiram {
			return BandModeROM
		}
	case page >= BasicStartPage && page <= BasicEndPage: // 0xa000 ~ 0xbfff
		if loram && hiram {
			return BandModeROM
		}
	case page >= CharStartPage && page <= CharEndPage: // 0xd000 ~ 0xdfff
		if charen && (loram || hiram) {
			return BandModeIO
		} else if !charen && (loram || hiram) {
			return BandModeROM
		}
	default:
		return BandModeRAM
	}
	return BandModeRAM
}
