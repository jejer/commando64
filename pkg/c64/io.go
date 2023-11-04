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
	keyboardIndex  map[uint32]uint8 // row: 0xf0, col: 0x0f
}

func NewIO(logger slog.Logger, c *Console) *IO {
	io := &IO{
		console: c,
		logger:  *logger.With("Component", "IO"),
	}
	for i := 0; i < 8; i++ {
		io.keyboardMetrix[i] = 0xff
	}

	io.initColor()
	io.initKeyboard()

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
		panic(1)
	}
	io.window = window

	// surface, err := window.GetSurface()
	surface, err := sdl.CreateRGBSurface(0, ScreenVisibleWidth, ScreenVisibleLines, 32, 0, 0, 0, 0)
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
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, ScreenVisibleWidth, ScreenVisibleLines)
	// texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		logger.Error("Texture create failed")
	}
	io.texture = texture

	return io
}

func (io *IO) initColor() {
	io.colors = [16]uint32{
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
	}
}

func (io *IO) initKeyboard() {
	io.keyboardIndex = make(map[uint32]uint8)
	io.keyboardIndex[sdl.SCANCODE_A] = 0x12
	io.keyboardIndex[sdl.SCANCODE_B] = 0x34
	io.keyboardIndex[sdl.SCANCODE_C] = 0x24
	io.keyboardIndex[sdl.SCANCODE_D] = 0x22
	io.keyboardIndex[sdl.SCANCODE_E] = 0x16
	io.keyboardIndex[sdl.SCANCODE_F] = 0x25
	io.keyboardIndex[sdl.SCANCODE_G] = 0x32
	io.keyboardIndex[sdl.SCANCODE_H] = 0x35
	io.keyboardIndex[sdl.SCANCODE_I] = 0x41
	io.keyboardIndex[sdl.SCANCODE_J] = 0x42
	io.keyboardIndex[sdl.SCANCODE_K] = 0x45
	io.keyboardIndex[sdl.SCANCODE_L] = 0x52
	io.keyboardIndex[sdl.SCANCODE_M] = 0x44
	io.keyboardIndex[sdl.SCANCODE_N] = 0x47
	io.keyboardIndex[sdl.SCANCODE_O] = 0x46
	io.keyboardIndex[sdl.SCANCODE_P] = 0x51
	io.keyboardIndex[sdl.SCANCODE_Q] = 0x76
	io.keyboardIndex[sdl.SCANCODE_R] = 0x21
	io.keyboardIndex[sdl.SCANCODE_S] = 0x15
	io.keyboardIndex[sdl.SCANCODE_T] = 0x26
	io.keyboardIndex[sdl.SCANCODE_U] = 0x36
	io.keyboardIndex[sdl.SCANCODE_V] = 0x37
	io.keyboardIndex[sdl.SCANCODE_W] = 0x11
	io.keyboardIndex[sdl.SCANCODE_X] = 0x27
	io.keyboardIndex[sdl.SCANCODE_Y] = 0x31
	io.keyboardIndex[sdl.SCANCODE_Z] = 0x14
	io.keyboardIndex[sdl.SCANCODE_0] = 0x43
	io.keyboardIndex[sdl.SCANCODE_1] = 0x70
	io.keyboardIndex[sdl.SCANCODE_2] = 0x73
	io.keyboardIndex[sdl.SCANCODE_3] = 0x10
	io.keyboardIndex[sdl.SCANCODE_4] = 0x13
	io.keyboardIndex[sdl.SCANCODE_5] = 0x20
	io.keyboardIndex[sdl.SCANCODE_6] = 0x23
	io.keyboardIndex[sdl.SCANCODE_7] = 0x30
	io.keyboardIndex[sdl.SCANCODE_8] = 0x33
	io.keyboardIndex[sdl.SCANCODE_9] = 0x40
	io.keyboardIndex[sdl.SCANCODE_F1] = 0x04
	io.keyboardIndex[sdl.SCANCODE_F3] = 0x05
	io.keyboardIndex[sdl.SCANCODE_F5] = 0x06
	io.keyboardIndex[sdl.SCANCODE_F7] = 0x03
	// other keys
	io.keyboardIndex[sdl.SCANCODE_RETURN] = 0x01
	io.keyboardIndex[sdl.SCANCODE_SPACE] = 0x74
	io.keyboardIndex[sdl.SCANCODE_LSHIFT] = 0x17
	io.keyboardIndex[sdl.SCANCODE_RSHIFT] = 0x64
	io.keyboardIndex[sdl.SCANCODE_COMMA] = 0x57
	io.keyboardIndex[sdl.SCANCODE_PERIOD] = 0x54
	io.keyboardIndex[sdl.SCANCODE_SLASH] = 0x67
	io.keyboardIndex[sdl.SCANCODE_SEMICOLON] = 0x62
	io.keyboardIndex[sdl.SCANCODE_EQUALS] = 0x65
	io.keyboardIndex[sdl.SCANCODE_BACKSPACE] = 0x00
	io.keyboardIndex[sdl.SCANCODE_MINUS] = 0x53
	// no physical keys, map to others
	io.keyboardIndex[sdl.SCANCODE_UP] = 0x56           // @
	io.keyboardIndex[sdl.SCANCODE_DOWN] = 0x55         // :
	io.keyboardIndex[sdl.SCANCODE_LEFT] = 0x61         // *
	io.keyboardIndex[sdl.SCANCODE_RIGHT] = 0x66        // ^
	io.keyboardIndex[sdl.SCANCODE_BACKSLASH] = 0x50    // +
	io.keyboardIndex[sdl.SCANCODE_RIGHTBRACKET] = 0x76 // C=  commando key
}

func (io *IO) Run() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				key := e.Keysym.Scancode
				if e.Repeat > 0 {
					io.logger.Info("key repeating", "key", key)
				} else {
					io.handleKey(uint32(e.Keysym.Scancode), e.State == sdl.PRESSED)
				}

			}
		}
		sdl.Delay(16)
	}
}

func (io *IO) handleKey(key uint32, pressed bool) {
	if pos, ok := io.keyboardIndex[key]; ok {
		row := pos >> 4
		col := pos & 0x0f
		if pressed {
			var mask uint8 = ^(1 << col)
			io.keyboardMetrix[row] &= mask
		} else {
			var mask uint8 = 1 << col
			io.keyboardMetrix[row] |= mask
		}
	}
}

func (io *IO) SetFramePixel(x int, y uint16, color uint8) {
	c := io.colors[color&0x0f]
	pixels := io.surface.Pixels()
	// index := int(y)*ScreenVisibleWidth*4 + x*4
	// pixels[index] = byte((c & 0xff000000) >> 24)
	// pixels[index+1] = byte((c & 0x00ff0000) >> 16)
	// pixels[index+2] = byte((c & 0x0000ff00) >> 8)
	// pixels[index+3] = byte((c & 0x000000ff))
	var p []uint32 = unsafe.Slice((*uint32)(unsafe.Pointer(&pixels[0])), len(pixels)/4)
	p[int(y)*ScreenVisibleWidth+x] = c
}

func (io *IO) RefreshScreen() {
	io.texture.Update(nil, unsafe.Pointer(&io.surface.Pixels()[0]), ScreenVisibleWidth*4)
	io.renderer.Clear()
	io.renderer.Copy(io.texture, nil, nil)
	io.renderer.Present()
}
