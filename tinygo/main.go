package main

import (
	"machine"
	"time"

	"github.com/soypat/tiny-lax/ascii"
	"tinygo.org/x/drivers/mpu6050"
)

var (
	// Pins 4 and 5 are default I2C1 pins
	bus = machine.I2C0
)

const (
	baud = 9600
)

func main() {
	// We'll use UART to transfer data through USB
	uart := machine.UART0
	// Configure UART for common baudrate
	err := uart.Configure(machine.UARTConfig{BaudRate: baud})
	if err != nil {
		panic(err)
	}
	println("start program")

	// configure I2C
	err = bus.Configure(machine.I2CConfig{Frequency: 100_000})
	if err != nil {
		panic(err)
	}

	imu := mpu6050.New(bus)
	imu.Configure()
	if !imu.Connected() {
		println("IMU not connected")
	}

	// Int32 buffer can be 11 digits long (including negative symbol). We separate
	var buffer [12*6 + 1]byte

	var data [6]int32
	for {
		for i := range buffer {
			buffer[i] = ' ' // set all characters in buffer to white space
		}
		buffer[len(buffer)-1] = '\n' // set last character to a newline

		data[0], data[1], data[2] = imu.ReadAcceleration()
		data[3], data[4], data[5] = imu.ReadRotation()
		// convert to radians
		data[3] /= 57
		data[4] /= 57
		data[5] /= 57
		for i := range data {
			ascii.PutInt32(data[i], buffer[:(i+1)*12])
		}
		uart.Write(buffer[:])
		time.Sleep(100 * time.Millisecond)
	}
}
