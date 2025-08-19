package main

import (
	CORE "PICO/CoreFiles"
	MODULES "PICO/Modules"
	DataStructures "PICO/SpecialDataStructs"
	"fmt"
	"machine"
)

// SETTINGS
var SAMPLE_RATE int = 44100
var BITS_PER_SAMPLE int = 12
var SIGNED bool = false

// GLOBALS DECLARATIONS
var Application *CORE.Program
var BuiltinLED machine.Pin
var MIC MODULES.Microphone
var Bluetooth *machine.UART
var BluetoothStatePin machine.Pin
var MicrophoneBuffer *DataStructures.Queue[byte]

// Functions
func CheckBluetoothConnection() bool {
	return BluetoothStatePin.Get()
}

func GetOneSample() []byte {
	if MicrophoneBuffer.Size() >= 2 {
		HighByte, _ := MicrophoneBuffer.Dequeue()
		LowByte, _ := MicrophoneBuffer.Dequeue()
		return []byte{HighByte, LowByte}
	} else {
		return nil
	}
}

// Entry Point
func main() {
	// APPLICATION INITIALIZATIONS
	Application = CORE.CreateApplication()

	// SETUP METHOD
	Application.LetSetup(
		func() {
			CORE.InitPeripherals()
			BuiltinLED = CORE.CreateIOPin(25, machine.PinOutput)

			MIC = machine.ADC{Pin: machine.GPIO28}
			MIC.Configure(machine.ADCConfig{
				Resolution: uint32(BITS_PER_SAMPLE),
				Reference:  3300,
				Samples:    4,
			})

			Bluetooth = machine.UART1
			Bluetooth.Configure(machine.UARTConfig{
				BaudRate: 115200,
				TX:       machine.GP1,
				RX:       machine.GP2,
			})
			BluetoothStatePin = CORE.CreateIOPin(15, machine.PinInput)
			MicrophoneBuffer = DataStructures.CreateQueue[byte](SAMPLE_RATE * 2)
		},
	)

	// LOOP METHOD
	Application.LetLoop(
		func() {
			duration := CORE.TimeIt(func() {
				for i := 0; i < SAMPLE_RATE; i++ {
					Sample := MIC.Get()
					HighByte := byte(Sample >> 8 & 0xFF)
					LowByte := byte(Sample & 0xFF)

					MicrophoneBuffer.Enqueue(HighByte)
					MicrophoneBuffer.Enqueue(LowByte)
				}
			})

			CORE.PrintLN(fmt.Sprintf("Time Required to take %d Sample : %d us", SAMPLE_RATE, duration))
			MicrophoneBuffer.Clear()
		},
	)

	// RUN PROGRAM
	Application.Run()
}
