package clock

import "time"

type Clock struct {
	pause bool
	VIC   chan bool
	CPU   chan int
	CIA1  chan bool
	CIA2  chan bool
}

func NewClock() *Clock {
	c := &Clock{}
	c.VIC = make(chan bool)
	c.CPU = make(chan int)
	c.CIA1 = make(chan bool)
	c.CIA2 = make(chan bool)
	return c
}

func (c *Clock) Pause(pause bool) {
	c.pause = pause
}

func (c *Clock) Run() {
	for {
		<-time.After(time.Duration(time.Microsecond)) // PAL cpu clock is 0.985MHz
		if !c.pause {
			c.CIA1 <- true
			c.CIA2 <- true
			// c.VIC <- true
		}
	}
}
