package c64

import (
	"log/slog"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type IO struct {
	console *Console
	logger  slog.Logger
	// video
	frame    [ScreenVisibleLines * ScreenVisibleWidth]uint32 // bitmap
	colors   [16]uint32
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	surface  *sdl.Surface

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
			0x000000ff,
			0xffffffff,
			0xab3126ff,
			0x66daffff,
			0xbb3fb8ff,
			0x55ce58ff,
			0x1d0e97ff,
			0xeaf57cff,
			0xb97418ff,
			0x785300ff,
			0xdd9387ff,
			0x5b5b5bff,
			0x8b8b8bff,
			0xb0f4acff,
			0xaa9defff,
			0xb8b8b8ff,
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
	io.keyboardIndex['\n'] = 0x01
	io.keyboardIndex[' '] = 0x74
	io.keyboardIndex['/'] = 0x67
	io.keyboardIndex['^'] = 0x66
	io.keyboardIndex['='] = 0x65
	io.keyboardIndex[';'] = 0x62
	io.keyboardIndex['*'] = 0x61
	io.keyboardIndex[','] = 0x57
	io.keyboardIndex['@'] = 0x56
	io.keyboardIndex[':'] = 0x55
	io.keyboardIndex['.'] = 0x54
	io.keyboardIndex['-'] = 0x53
	io.keyboardIndex['+'] = 0x50

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		logger.Error("SDL init failed")
	}
	window, err := sdl.CreateWindow(
		"commando64",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		ScreenVisibleWidth*2,
		ScreenVisibleLines*2,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		logger.Error("Window create failed")
	}
	io.window = window

	surface, err := window.GetSurface()
	if err != nil {
		logger.Error("Surface create failed")
	}
	io.surface = surface
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		logger.Error("Renderer create failed")
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()
	io.renderer = renderer
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, ScreenVisibleWidth, ScreenVisibleLines)
	if err != nil {
		logger.Error("Texture create failed")
	}
	io.texture = texture

	return io
}

func (io *IO) Run() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}
		sdl.Delay(16)
	}
}

func (io *IO) SetFramePixel(x int, y uint16, color uint8) {
	c := io.colors[color&0x0f]
	pixels := io.surface.Pixels()
	pixels[int(y)*ScreenVisibleWidth+x] = byte((c & 0xff000000) >> 24)
	pixels[int(y)*ScreenVisibleWidth+x+1] = byte((c & 0xff000000) >> 16)
	pixels[int(y)*ScreenVisibleWidth+x+2] = byte((c & 0xff000000) >> 8)
	pixels[int(y)*ScreenVisibleWidth+x+3] = byte((c & 0xff000000))
	io.frame[int(y)*ScreenVisibleWidth+x] = io.colors[color&0x0f]
}

func (io *IO) RefreshScreen() {
	io.texture.Update(nil, unsafe.Pointer(&io.surface.Pixels()[0]), ScreenVisibleWidth*4)
	// io.renderer.Clear()
	// io.renderer.Copy(io.texture, nil, nil)
	// io.renderer.Present()
	io.window.UpdateSurface()
}
