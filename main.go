package main

import (
	"fmt"
)

func execOrDie(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}

func main() {
	device := newDevice()

	execOrDie(func() error { return device.Init(1, 0x5d) })
	execOrDie(device.Active)
	defer device.Deactive()

	pressure, _ := device.ReadPressure()
	temperature, _ := device.ReadTemperature()

	fmt.Printf("%.2f\t%.2f\n", pressure, temperature)
}
