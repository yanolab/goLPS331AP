package main

import (
	"fmt"
	"time"
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

	for i := 0; i < 1; i++ {
		pressure, _ := device.ReadPressure()
		temp, _ := device.ReadTemp()

		fmt.Printf("%.2f\t%.2f\n", pressure, temp)

		time.Sleep(1000 * 1000 * 1000)
	}
}
