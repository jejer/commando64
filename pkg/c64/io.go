package c64

import "log/slog"

type IO struct {
	console *Console
	logger  slog.Logger
	// video
	frame  [ScreenVisibleLines * ScreenVisibleWidth]uint32 // bitmap
	colors [16]uint32

	// keyboard
	// 		PB7	PB6	PB5	PB4	PB3	PB2	PB1	PB0
	// PA7	STOP	Q	C=	SPACE	2	CTRL	<-	1
	// PA6	/	^	=	RSHIFT	HOME	;	*	£
	// PA5	,	@	:	.	-	L	P	+
	// PA4	N	O	K	M	0	J	I	9
	// PA3	V	U	H	B	8	G	Y	7
	// PA2	X	T	F	C	6	D	R	5
	// PA1	LSHIFT	E	S	Z	4	A	W	3
	// PA0	CRSR DN	F5	F3	F1	F7	CRSR RT	RETURN	DELETE
	keyboardMetrix [8]uint8
	keyboardIndex  map[uint8]uint8 // row: 0xf0, col: 0x0f
}

func NewIO(logger slog.Logger, c *Console) *IO {
	io := &IO{
		console: c,
		logger:  *logger.With("Component", "IO"),
		colors: [16]uint32{
			0xff000000,
			0xffffffff,
			0xffab3126,
			0xff66daff,
			0xffbb3fb8,
			0xff55ce58,
			0xff1d0e97,
			0xffeaf57c,
			0xffb97418,
			0xff785300,
			0xffdd9387,
			0xff5b5b5b,
			0xff8b8b8b,
			0xffb0f4ac,
			0xffaa9def,
			0xffb8b8b8,
		},
	}
	for i := 0; i < 8; i++ {
		io.keyboardMetrix[i] = 0xff
	}
	io.keyboardIndex = make(map[uint8]uint8)
	io.keyboardIndex['A'] = 0x12
	io.keyboardIndex['B'] = 0x34
	io.keyboardIndex['C'] = 0x24
	io.keyboardIndex['D'] = 0x22
	io.keyboardIndex['E'] = 0x16
	io.keyboardIndex['F'] = 0x25
	io.keyboardIndex['G'] = 0x32
	io.keyboardIndex['H'] = 0x35
	io.keyboardIndex['I'] = 0x41
	io.keyboardIndex['J'] = 0x42
	io.keyboardIndex['K'] = 0x45
	io.keyboardIndex['L'] = 0x52
	io.keyboardIndex['M'] = 0x44
	io.keyboardIndex['N'] = 0x47
	io.keyboardIndex['O'] = 0x46
	io.keyboardIndex['P'] = 0x51
	io.keyboardIndex['Q'] = 0x76
	io.keyboardIndex['R'] = 0x21
	io.keyboardIndex['S'] = 0x15
	io.keyboardIndex['T'] = 0x26
	io.keyboardIndex['U'] = 0x36
	io.keyboardIndex['V'] = 0x37
	io.keyboardIndex['W'] = 0x11
	io.keyboardIndex['X'] = 0x27
	io.keyboardIndex['Y'] = 0x31
	io.keyboardIndex['Z'] = 0x14
	io.keyboardIndex['0'] = 0x43
	io.keyboardIndex['1'] = 0x70
	io.keyboardIndex['2'] = 0x73
	io.keyboardIndex['3'] = 0x10
	io.keyboardIndex['4'] = 0x13
	io.keyboardIndex['5'] = 0x20
	io.keyboardIndex['6'] = 0x23
	io.keyboardIndex['7'] = 0x30
	io.keyboardIndex['8'] = 0x33
	io.keyboardIndex['9'] = 0x40
	return io
}

func (io *IO) SetFramePixel(x, y int, color uint8) {
	io.frame[y*ScreenVisibleWidth+x] = io.colors[color&0x0f]
}
