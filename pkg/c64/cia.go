package c64

type CIA1 struct {
	console *Console

	// https://www.c64-wiki.com/wiki/CIA
	// $DC00 Data Port A, keyboard matrix columns
	dataPortA uint8
	// $DC01 Data Port B, keyboard matrix rows
	dataPortB uint8
	// $DC02 Data Direction Port A, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortADir uint8
	// $DC03 Data Direction Port B, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortBDir uint8
	// $DC04 $DC05 TimerA
	timerA [2]uint8
	// $DC06 $DC07 TimerB
	timerB [2]uint8
	// $DC08 ~ $DC0B Real Time Clock, 0.1s, 1s, 1m, 1h
	rtc [4]uint8
	// $DC0C Serial shift register
	sdr uint8
	// $DC0D Interrupt Control and status
	irq uint8
	// $DC0E Control Timer A
	timerAControl uint8
	// $DC0F Control Timer B
	timerBControl uint8
}

func (cia1 *CIA1) Write(addr uint16, v uint8) {

}
func (cia1 *CIA1) Read(addr uint16) uint8 {
	return 0
}

type CIA2 struct {
	console *Console

	// https://www.c64-wiki.com/wiki/CIA
	// $DD00 Data Port A, VIC bank selection and serial bus
	dataPortA uint8
	// $DD01 Data Port B, RS232 related
	dataPortB uint8
	// $DD02 Data Direction Port A, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortADir uint8
	// $DD03 Data Direction Port B, Bit X: 0=Input (read only), 1=Output (read and write)
	dataPortBDir uint8
	// $DD04 $DD05 TimerA
	timerA [2]uint8
	// $DD06 $DD07 TimerB
	timerB [2]uint8
	// $DD08 ~ $DD0B Real Time Clock, 0.1s, 1s, 1m, 1h
	rtc [4]uint8
	// $DD0C Serial shift register
	sdr uint8
	// $DD0D Interrupt Control and status, CIA2 is connected to the NMI-Line.
	irq uint8
	// $DD0E Control Timer A
	timerAControl uint8
	// $DD0F Control Timer B
	timerBControl uint8
}

func (cia2 *CIA2) Write(addr uint16, v uint8) {

}
func (cia2 *CIA2) Read(addr uint16) uint8 {
	return 0
}
