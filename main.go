/*
Somehow, I broke the GO Error checker and all of the "IDE" stuff is gone.
Ain't my luck is great. Despite all of this, TinyGO is able to compile this project.
So I am not re-"making" (IDK What is the word I should write there) my project.
*/

package main

import (
	CORE "PI_PICO/CoreFiles"
	MODULES "PI_PICO/Modules"
	"machine"
	"time"
	"fmt"
)

// SETTINGS
var SAMPLE_RATE int = 8000
var BITS_PER_SAMPLE int = 16
var SIGNED bool = false
var BIG_ENDIAN bool = true
var DEBUG_PRINT bool = true

// GLOBALS DECLARATIONS
var Application *CORE.Program
var BuiltinLED machine.Pin
var MIC MODULES.Microphone
var Bluetooth *machine.UART
var BluetoothStatePin machine.Pin
var ActiveButton machine.Pin
var Sample uint16

// FUNCTIONS
func CheckBluetoothConnection() bool {
	return BluetoothStatePin.Get()
}

// ENTRY POINT
func main() {
	// APPLICATION INITIALIZATIONS
	Application = CORE.CreateApplication()

	// SETUP METHOD
	Application.LetSetup(
		func() {
			CORE.InitPeripherals()
			BuiltinLED = CORE.CreateIOPin(25, machine.PinOutput)
			ActiveButton = CORE.CreateIOPin(16, machine.PinInput)

			MIC = machine.ADC{Pin: machine.ADC2}
			MIC.Configure(machine.ADCConfig{
				Resolution: uint32(BITS_PER_SAMPLE),
				Reference:  3300,
				Samples:    4,
			})

			Bluetooth = machine.UART1
			Bluetooth.Configure(machine.UARTConfig{
				BaudRate: 115200,
				TX:       machine.GP4,
				RX:       machine.GP5,
			})
			BluetoothStatePin = CORE.CreateIOPin(13, machine.PinInput)
		},
	)

	// LOOP METHOD
	Application.LetLoop(
		func() {
			if ActiveButton.Get() && CheckBluetoothConnection() {
				BuiltinLED.Low()

				StartTime := time.Now()

				duration := CORE.TimeIt(func() {
					Sample = MIC.Get()
					CORE.PrintLN(fmt.Sprintf("Sample - %d", Sample))

					if BIG_ENDIAN {
						Bluetooth.Write([]byte{ byte(Sample >> 8 & 0xFF), byte(Sample & 0xFF) }) 
					} else {
						Bluetooth.Write([]byte{ byte(Sample & 0xFF), byte(Sample >> 8 & 0xFF) })
					}
				})
				delay_period := time.Duration(1000000 / SAMPLE_RATE) * time.Microsecond - duration
				if delay_period > 0 {
					CORE.Delay( time.Duration(delay_period) )
				}

				EndTime := time.Since(StartTime)

				if DEBUG_PRINT {
					CORE.PrintLN(fmt.Sprintf("Send Cycle Time : %d us", EndTime.Microseconds())) 
				}
			} else {
				CORE.PrintLN("Not Active")
				BuiltinLED.High()
			}
		},
	)

	// RUN PROGRAM
	Application.Run()
}
