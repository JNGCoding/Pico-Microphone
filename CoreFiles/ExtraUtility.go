package CORE

import (
	"machine"
	"time"
)

func InitPeripherals() {
	machine.InitADC()
	machine.InitSerial()
}

func PrintLN(text string) {
	machine.USBCDC.Write([]byte(text + "\n"))
}

func Print(text string) {
	machine.USBCDC.Write([]byte(text))
}

func WriteBytes(data []byte) {
	machine.USBCDC.Write(data)
}

func CreateIOPin(num int16, config machine.PinMode) machine.Pin {
	result := machine.Pin(num)
	result.Configure(
		machine.PinConfig{
			Mode: config,
		},
	)

	return result
}

func Delay(n time.Duration) {
	time.Sleep(n * time.Millisecond)
}

func TimeIt(fn func()) time.Duration {
	start_time := time.Now()
	fn()
	return time.Since(start_time)
}
