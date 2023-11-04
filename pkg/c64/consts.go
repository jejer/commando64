package c64

const (
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

	// roms
	BasicRomAddr  uint16 = 0xa000
	KernalRomAddr uint16 = 0xe000
	CharsRomAddr  uint16 = 0xd000
)
