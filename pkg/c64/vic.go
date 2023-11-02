package c64

import "log/slog"

// VIC-II
// https://web.archive.org/web/20230521191139/https://dustlayer.com/vic-ii/2013/4/22/when-visibility-matters

// VIC-II has 14 address lines = 16 kB address space
// registers $D000 ~ $D02E
// http://unusedino.de/ec64/technical/aay/c64/vicmain.htm
type VICII struct {
	console *Console
	logger  slog.Logger

	// https://www.c64-wiki.com/wiki/Page_208-211
	// $00~$10 sprite position
	spritePos  [16]uint8
	spriteMSBx uint8 // the 5th bit of x position
	// $11 control register 1
	cr1 uint8
	// $12 ruster position
	rasterPos uint8
	// $13 $14 Light Pen position
	lightpenPos [2]uint8
	// $15 sprite enabled
	spriteEnabled uint8
	// $16 control register 2
	cr2 uint8
	// $17 sprite y expansion
	spriteExpY uint8
	// $18 memory pointers
	memPointers uint8
	// $19 interrupt register
	interruptState uint8
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
}

func NewVICII(c *Console, logger slog.Logger) *VICII {
	vic := &VICII{console: c}
	vic.logger = *logger.With("Component", "VICII")
	return vic
}

func (vic *VICII) Write(addr uint16, v uint8) {

}
func (vic *VICII) Read(addr uint16) uint8 {
	addr &= 0x00ff
	switch addr {
	// case 0x12:
	// 	vic.rasterPos++
	// 	return vic.rasterPos
	}
	return 0
}
