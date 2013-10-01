package main

import (
	"bitbucket.org/gmcbay/i2c"
	"errors"
)

type LPS331AP struct {
	bus    *i2c.I2CBus
	addr   byte
	active bool
}

func newDevice() *LPS331AP {
	return &LPS331AP{}
}

func (l *LPS331AP) Init(busNumber byte, addr byte) error {
	var err error
	l.bus, err = i2c.Bus(busNumber)
	l.addr = addr

	return err
}

func (l *LPS331AP) Read(reg byte) (byte, error) {
	buf, err := l.bus.ReadByteBlock(l.addr, reg, 1)
	if err != nil {
		return 0, err
	}

	return buf[0], nil
}

func (l *LPS331AP) ReadPressure() (float32, error) {
	buf := make([]byte, 3)

	for idx := 0x28; idx <= 0x2a; idx++ {
		var err error
		buf[idx-0x28], err = l.Read(byte(idx))
		if err != nil {
			return 0, err
		}
	}

	return float32(int(buf[2])<<16|int(buf[1])<<8|int(buf[0])) / 4096.0, nil
}

func (l *LPS331AP) ReadTemp() (float32, error) {
	buf := make([]byte, 2)

	for idx := 0x2b; idx <= 0x2c; idx++ {
		var err error
		buf[idx-0x2b], err = l.Read(byte(idx))
		if err != nil {
			return 0, err
		}
	}

	return 42.5 + float32(^(int16(buf[1])<<8|int16(buf[0]))+1)*-1.0/480.0, nil
}

func (l *LPS331AP) Active() error {
	id, err := l.Read(0x0f)
	if err != nil {
		return err
	}
	if id != 0xbb {
		return errors.New("Invalid device.")
	}

	if err := l.bus.WriteByte(l.addr, 0x20, 0x90); err != nil {
		return err
	}

	l.active = true

	return nil
}

func (l *LPS331AP) Deactive() error {
	if !l.active {
		return nil
	}

	if err := l.bus.WriteByte(l.addr, 0x20, 0x0); err != nil {
		return err
	}

	l.active = false

	return nil
}
