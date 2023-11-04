package vic

import (
	"fmt"
	"log/slog"

	"github.com/jejer/commando64/pkg/c64"
	"github.com/jejer/commando64/pkg/c64/clock"
)

// VIC-II
// https://web.archive.org/web/20230521191139/https://dustlayer.com/vic-ii/2013/4/22/when-visibility-matters

// VIC-II has 14 address lines = 16 kB address space
// registers $D000 ~ $D02E
// http://unusedino.de/ec64/technical/aay/c64/vicmain.htm
type GraphicMode int

const (
	// VIC
	StdCharMode          GraphicMode = iota // ECM0 BMM0 MCM0
	MultiColorCharMode                      // ECM0 BMM0 MCM1
	StdBitmapMode                           // ECM0 BMM1 MCM0
	MultiColorBitmapMode                    // ECM0 BMM1 MCM1
	ExtBGColorMode                          // ECM1 BMM0 MCM0
	InvalidMode

	// screen constants
	// https://dustlayer.com/vic-ii/2013/4/25/vic-ii-for-beginners-beyond-the-screen-rasters-cycle
	ScreenLines            = 312
	ScreenWidth            = 504
	ScreenVisibleLines     = 284
	ScreenVisibleWidth     = 403
	ScreenFirstVisibleLine = 14 // start from 0
	ScreenLastVisibleLine  = 298
	ScreenFirstTextLine    = 56
	ScreenLastTextLine     = 256
	ScreenFirstTextCol     = 42
	ScreenTextLines        = 200
	ScreenTextWidth        = 320
	ScreenTextPerLine      = 40
	LineCycles             = 63
	BadLineCycles          = 23

	ColorRamStartPage uint16 = 0xd800
)

type VICII struct {
	logger       slog.Logger
	clock        *clock.Clock
	cycle        int8
	cpuCycle     int8
	mem          c64.MemoryBus
	irqCh        chan<- bool
	peripheralIO c64.PeripheralIO

	// https://www.c64-wiki.com/wiki/Page_208-211
	// $00~$10 sprite position
	spritePos  [16]uint8
	spriteMSBx uint8 // the 5th bit of x position
	// $11 control register 1
	control1 uint8
	// $12 ruster position
	rasterPos uint8
	// $13 $14 Light Pen position
	lightpenPos [2]uint8
	// $15 sprite enabled
	spriteEnabled uint8
	// $16 control register 2
	control2 uint8
	// $17 sprite y expansion
	spriteExpY uint8
	// $18 memory pointers
	memPointers uint8
	// $19 interrupt register
	interruptStatus uint8
	// $1a interrupt enabled
	interruptEnabled uint8
	// $1b sprite data priority
	spriteDataPriority uint8
	// $1c sprite multicolor
	spriteMulticolor uint8
	// $1d sprite x expansion
	spriteExpX uint8
	// $1e sprite-sprite collision
	spriteSpriteCollision uint8
	// $1f sprite-data collision
	spriteDataCollision uint8
	// $20 border color
	colorBorder uint8
	// $21~$24 background colors
	colorBackground [4]uint8
	// $25 $26 sprite multicolor
	colorSpriteMulti [2]uint8
	// $27~$2e sprite colors
	colorSprite [8]uint8

	// status
	mode            GraphicMode
	charMemOffset   uint16 // offsets by memory pointers
	screenMemOffset uint16
	bitmapMemOffset uint16

	rasterIrqRequest uint16
}

func NewVICII(logger slog.Logger, clock *clock.Clock, m c64.MemoryBus, ch chan<- bool, io c64.PeripheralIO) *VICII {
	vic := &VICII{mem: m, peripheralIO: io, irqCh: ch, clock: clock}
	vic.logger = *logger.With("Component", "VICII")
	vic.cycle = 1
	vic.cpuCycle = 1
	return vic
}

func (vic *VICII) Write(addr uint16, v uint8) {
	switch add := addr & 0x00ff; {
	case add >= 0x00 && add <= 0x0f:
		vic.spritePos[add] = v
	case add == 0x10:
		vic.spriteMSBx = v
	case add == 0x11:
		vic.control1 = v & 0x7f // 7th bit is the 8th bit of raster counter
		vic.rasterIrqRequest &= 0x00ff
		vic.rasterIrqRequest |= ((uint16(v) & 0x0080) << 1)
		vic.setGraphicMode()
	case add == 0x12:
		vic.rasterIrqRequest &= 0x0100
		vic.rasterIrqRequest |= uint16(v)
	case add == 0x13: // x
		vic.lightpenPos[0] = v
	case add == 0x14: // y
		vic.lightpenPos[1] = v
	case add == 0x15:
		vic.spriteEnabled = v
	case add == 0x16:
		vic.control2 = v
		vic.setGraphicMode()
	case add == 0x17:
		vic.spriteExpY = v
	case add == 0x18:
		vic.memPointers = v | 1
		// bits ----xxx-
		vic.charMemOffset = (uint16(v) & 0x0e) << 10
		// bits xxxx----
		vic.screenMemOffset = (uint16(v) & 0xf0) << 6
		// bit  ----x---
		vic.bitmapMemOffset = (uint16(v) & 0x08) << 10
	case add == 0x19: // irq acknowledge
		vic.interruptStatus &= ^(v & 0x0f)
	case add == 0x1a:
		vic.interruptEnabled = v
	case add == 0x1b:
		vic.spriteDataPriority = v
	case add == 0x1c:
		vic.spriteMulticolor = v
	case add == 0x1d:
		vic.spriteExpX = v
	case add == 0x1e || add == 0x1f: // sprite collision
	case add == 0x20:
		vic.colorBorder = v
	case add >= 0x21 && add <= 0x24:
		vic.colorBackground[add-0x21] = v
	case add == 0x25 || add == 0x26:
		vic.colorSpriteMulti[add-0x25] = v
	case add >= 0x27 && add <= 0x2e:
		vic.colorSprite[add-0x27] = v
	default:
		vic.logger.Warn("VIC register write", "address", fmt.Sprintf("0x%x", addr))
	}
}
func (vic *VICII) Read(addr uint16) uint8 {
	switch add := addr & 0x00ff; {
	case add >= 0x00 && add <= 0x0f:
		return vic.spritePos[add]
	case add == 0x10:
		return vic.spriteMSBx
	case add == 0x11:
		return vic.control1
	case add == 0x12:
		return vic.rasterPos
	case add == 0x13: // x
		return vic.lightpenPos[0]
	case add == 0x14: // y
		return vic.lightpenPos[1]
	case add == 0x15:
		return vic.spriteEnabled
	case add == 0x16:
		return vic.control2
	case add == 0x17:
		return vic.spriteExpY
	case add == 0x18:
		return vic.memPointers
	case add == 0x19: // irq status
		v := vic.interruptStatus & 0x0f
		if v != 0 {
			v |= 0x80 // IRQ bit
		}
		v |= 0x70 // non-connected bits (always set)
		return v
	case add == 0x1a:
		return vic.interruptEnabled | 0xf0
	case add == 0x1b:
		return vic.spriteDataPriority
	case add == 0x1c:
		return vic.spriteMulticolor
	case add == 0x1d:
		return vic.spriteExpX
	case add == 0x1e || add == 0x1f: // sprite collision
	case add == 0x20:
		return vic.colorBorder
	case add >= 0x21 && add <= 0x24:
		return vic.colorBackground[add-0x21]
	case add == 0x25 || add == 0x26:
		return vic.colorSpriteMulti[add-0x25]
	case add >= 0x27 && add <= 0x2e:
		return vic.colorSprite[add-0x27]
	default:
		vic.logger.Warn("VIC register read", "address", fmt.Sprintf("0x%x", addr))
	}
	return 0
}

func (vic *VICII) Run() {
	for {
		<-vic.clock.VIC
		vic.cpuCycle--
		if vic.cpuCycle >= 0 {
			vic.clock.CPU <- true
		}
		vic.cycle--
		if vic.cycle == 0 {
			vic.step()
		}
	}
}

// a step is a raster line
func (vic *VICII) step() {
	if vic.interruptStatus&0x80 != 0 {
		// interrupts are not handled by CPU
		go func() { vic.irqCh <- false }()
	}

	var line uint16 = uint16(vic.rasterPos) | (uint16(vic.control1&0x0080) << 1)

	if vic.interruptEnabled&0x01 != 0 && line == vic.rasterIrqRequest {
		// check and trigger raster line irq
		vic.interruptStatus |= 0x01
		go func() { vic.irqCh <- false }()
	}

	// draw line
	if line >= ScreenFirstVisibleLine && line < ScreenLastVisibleLine {
		y := line - ScreenFirstVisibleLine
		for x := 0; x < ScreenVisibleWidth; x++ {
			vic.peripheralIO.SetFramePixel(x, y, vic.colorBorder)
		}
		switch vic.mode {
		case StdCharMode:
			vic.drawCharRasterLine(line, y)
		default:
			vic.logger.Error("VIC mod not implemented", "mode", vic.mode)
		}
	}

	// update next cycle
	vic.cycle = LineCycles
	if vic.isBadLine(line) {
		vic.cpuCycle = BadLineCycles
	} else {
		vic.cpuCycle = LineCycles
	}

	// update raster
	line++
	if line == ScreenLines {
		line = 0
		vic.peripheralIO.RefreshScreen()
	}
	vic.rasterPos = uint8(line & 0x00ff)
	vic.control1 &= 0x7f
	vic.control1 |= uint8((line >> 1) & 0x80)
}

func (vic *VICII) setGraphicMode() {
	mode := (vic.control1 & 0x60) >> 4 // get ICM and BMM bit
	mode |= (vic.control2 & 0x10) >> 4 // get MCM bit
	if GraphicMode(mode) >= InvalidMode {
		vic.logger.Error("SetGraphicMode: invalid mode", "mode", mode)
	}
	vic.mode = GraphicMode(mode)
	vic.logger.Info("SetGraphicMode", "mode", mode)
}

func (vic *VICII) drawCharRasterLine(line, y uint16) {
	if line < ScreenFirstTextLine || line >= ScreenLastTextLine || (vic.control1<<4) == 0 {
		return
	}

	// text background
	for x := 0; x < ScreenTextWidth; x++ {
		vic.peripheralIO.SetFramePixel(x+ScreenFirstTextCol, y, vic.colorBackground[0])
	}

	// text dots in this line
	for col := 0; col < ScreenTextPerLine; col++ {
		row := (line - ScreenFirstTextLine) / 8
		char := vic.getScreenChar(row, uint16(col))
		color := vic.getCharColor(row, uint16(col))
		data := vic.getCharData(char, (line-ScreenFirstTextLine)%8)
		for i := 0; i < 8; i++ {
			if data&(1<<i) != 0 {
				x := ScreenFirstTextCol + (col * 8) + 8 - i
				vic.peripheralIO.SetFramePixel(x, y, color)
			}
		}
	}
}

func (vic *VICII) getScreenChar(row, col uint16) uint8 {
	addr := vic.screenMemOffset + row*ScreenTextPerLine + col
	return vic.mem.VicRead(addr)
}

func (vic *VICII) getCharData(char uint8, row uint16) uint8 {
	addr := vic.charMemOffset + (uint16(char) * 8) + row
	return vic.mem.VicRead(addr)
}

func (vic *VICII) getCharColor(row, col uint16) uint8 {
	addr := ColorRamStartPage + row*ScreenTextPerLine + col
	return vic.mem.Read(addr)
}

// According to Christian Bauer's paper:
//
// A Bad Line Condition is given at any arbitrary clock cycle,
// if at the negative edge of 0 at the beginning of the cycle
// RASTER >= $30 and RASTER <= $f7 and the lower three bits
// of RASTER are equal to YSCROLL and if the DEN bit was set
// during an arbitrary cycle of raster line $30.
func (vic *VICII) isBadLine(line uint16) bool {
	return (line >= 0x0030 && line <= 0x00f7 && ((line & 0x0007) == (uint16(vic.control1) & 0x0007)))
}
