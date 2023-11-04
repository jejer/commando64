package peripheral

import (
	"log/slog"
	"unsafe"

	"github.com/jejer/commando64/pkg/c64"
	"github.com/veandco/go-sdl2/sdl"
)

type PeripheralSDL struct {
	logger slog.Logger

	colors         [16]uint32
	keyboardMetrix [8]uint8
	keyboardIndex  map[uint32]uint8 // row: 0xf0, col: 0x0f

	renderer *sdl.Renderer
	texture  *sdl.Texture
	surface  *sdl.Surface
}

func NewPeripheralSDL(logger slog.Logger) *PeripheralSDL {
	sdl := &PeripheralSDL{
		logger: *logger.With("Component", "PeripheralSDL"),
	}
	return sdl
}

func (p *PeripheralSDL) Init() {
	p.initColor()
	p.initKeyboard()
	p.initVideo()
}

func (p *PeripheralSDL) EventLoop() {
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
					p.logger.Info("key repeating", "key", key)
				} else {
					p.handleKey(uint32(e.Keysym.Scancode), e.State == sdl.PRESSED)
				}

			}
		}
	}
}

func (p *PeripheralSDL) handleKey(key uint32, pressed bool) {
	if pos, ok := p.keyboardIndex[key]; ok {
		row := pos >> 4
		col := pos & 0x0f
		if pressed {
			var mask uint8 = ^(1 << col)
			p.keyboardMetrix[row] &= mask
		} else {
			var mask uint8 = 1 << col
			p.keyboardMetrix[row] |= mask
		}
	}
}

func (p *PeripheralSDL) ReadKeyboardMatrix(row uint8) uint8 {
	return p.keyboardMetrix[row]
}

func (p *PeripheralSDL) SetFramePixel(x int, y uint16, color uint8) {
	c := p.colors[color&0x0f]
	pixels := p.surface.Pixels()
	var d []uint32 = unsafe.Slice((*uint32)(unsafe.Pointer(&pixels[0])), len(pixels)/4)
	d[int(y)*c64.ScreenVisibleWidth+x] = c
}
func (p *PeripheralSDL) RefreshScreen() {
	p.texture.Update(nil, unsafe.Pointer(&p.surface.Pixels()[0]), c64.ScreenVisibleWidth*4)
	p.renderer.Clear()
	p.renderer.Copy(p.texture, nil, nil)
	p.renderer.Present()
}

func (p *PeripheralSDL) initKeyboard() {
	for i := 0; i < 8; i++ {
		p.keyboardMetrix[i] = 0xff
	}

	p.keyboardIndex = make(map[uint32]uint8)
	p.keyboardIndex[sdl.SCANCODE_A] = 0x12
	p.keyboardIndex[sdl.SCANCODE_B] = 0x34
	p.keyboardIndex[sdl.SCANCODE_C] = 0x24
	p.keyboardIndex[sdl.SCANCODE_D] = 0x22
	p.keyboardIndex[sdl.SCANCODE_E] = 0x16
	p.keyboardIndex[sdl.SCANCODE_F] = 0x25
	p.keyboardIndex[sdl.SCANCODE_G] = 0x32
	p.keyboardIndex[sdl.SCANCODE_H] = 0x35
	p.keyboardIndex[sdl.SCANCODE_I] = 0x41
	p.keyboardIndex[sdl.SCANCODE_J] = 0x42
	p.keyboardIndex[sdl.SCANCODE_K] = 0x45
	p.keyboardIndex[sdl.SCANCODE_L] = 0x52
	p.keyboardIndex[sdl.SCANCODE_M] = 0x44
	p.keyboardIndex[sdl.SCANCODE_N] = 0x47
	p.keyboardIndex[sdl.SCANCODE_O] = 0x46
	p.keyboardIndex[sdl.SCANCODE_P] = 0x51
	p.keyboardIndex[sdl.SCANCODE_Q] = 0x76
	p.keyboardIndex[sdl.SCANCODE_R] = 0x21
	p.keyboardIndex[sdl.SCANCODE_S] = 0x15
	p.keyboardIndex[sdl.SCANCODE_T] = 0x26
	p.keyboardIndex[sdl.SCANCODE_U] = 0x36
	p.keyboardIndex[sdl.SCANCODE_V] = 0x37
	p.keyboardIndex[sdl.SCANCODE_W] = 0x11
	p.keyboardIndex[sdl.SCANCODE_X] = 0x27
	p.keyboardIndex[sdl.SCANCODE_Y] = 0x31
	p.keyboardIndex[sdl.SCANCODE_Z] = 0x14
	p.keyboardIndex[sdl.SCANCODE_0] = 0x43
	p.keyboardIndex[sdl.SCANCODE_1] = 0x70
	p.keyboardIndex[sdl.SCANCODE_2] = 0x73
	p.keyboardIndex[sdl.SCANCODE_3] = 0x10
	p.keyboardIndex[sdl.SCANCODE_4] = 0x13
	p.keyboardIndex[sdl.SCANCODE_5] = 0x20
	p.keyboardIndex[sdl.SCANCODE_6] = 0x23
	p.keyboardIndex[sdl.SCANCODE_7] = 0x30
	p.keyboardIndex[sdl.SCANCODE_8] = 0x33
	p.keyboardIndex[sdl.SCANCODE_9] = 0x40
	p.keyboardIndex[sdl.SCANCODE_F1] = 0x04
	p.keyboardIndex[sdl.SCANCODE_F3] = 0x05
	p.keyboardIndex[sdl.SCANCODE_F5] = 0x06
	p.keyboardIndex[sdl.SCANCODE_F7] = 0x03
	// other keys
	p.keyboardIndex[sdl.SCANCODE_RETURN] = 0x01
	p.keyboardIndex[sdl.SCANCODE_SPACE] = 0x74
	p.keyboardIndex[sdl.SCANCODE_LSHIFT] = 0x17
	p.keyboardIndex[sdl.SCANCODE_RSHIFT] = 0x64
	p.keyboardIndex[sdl.SCANCODE_COMMA] = 0x57
	p.keyboardIndex[sdl.SCANCODE_PERIOD] = 0x54
	p.keyboardIndex[sdl.SCANCODE_SLASH] = 0x67
	p.keyboardIndex[sdl.SCANCODE_SEMICOLON] = 0x62
	p.keyboardIndex[sdl.SCANCODE_EQUALS] = 0x65
	p.keyboardIndex[sdl.SCANCODE_BACKSPACE] = 0x00
	p.keyboardIndex[sdl.SCANCODE_MINUS] = 0x53
	// no physical keys, map to others
	p.keyboardIndex[sdl.SCANCODE_UP] = 0x56           // @
	p.keyboardIndex[sdl.SCANCODE_DOWN] = 0x55         // :
	p.keyboardIndex[sdl.SCANCODE_LEFT] = 0x61         // *
	p.keyboardIndex[sdl.SCANCODE_RIGHT] = 0x66        // ^
	p.keyboardIndex[sdl.SCANCODE_BACKSLASH] = 0x50    // +
	p.keyboardIndex[sdl.SCANCODE_RIGHTBRACKET] = 0x76 // C=  commando key
}

func (p *PeripheralSDL) initColor() {
	p.colors = [16]uint32{
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

func (p *PeripheralSDL) initVideo() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		p.logger.Error("SDL init failed")
	}
	window, err := sdl.CreateWindow(
		"commando64",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		c64.ScreenVisibleWidth*2,
		c64.ScreenVisibleLines*2,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		p.logger.Error("Window create failed")
		panic(1)
	}

	// surface, err := window.GetSurface()
	surface, err := sdl.CreateRGBSurface(0, c64.ScreenVisibleWidth, c64.ScreenVisibleLines, 32, 0, 0, 0, 0)
	if err != nil {
		p.logger.Error("Surface create failed")
	}
	p.surface = surface
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		p.logger.Error("Renderer create failed")
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()
	p.renderer = renderer
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, c64.ScreenVisibleWidth, c64.ScreenVisibleLines)
	// texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		p.logger.Error("Texture create failed")
	}
	p.texture = texture
}
