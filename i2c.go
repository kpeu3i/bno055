package bno055

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

const (
	i2cSlave = 0x0703
)

type i2c struct {
	mu sync.Mutex
	rc *os.File
}

func newI2C(addr uint8, bus int) (*i2c, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	err = ioctl(f.Fd(), i2cSlave, uintptr(addr))
	if err != nil {
		return nil, err
	}

	return &i2c{rc: f}, nil
}

func (bus *i2c) Read(reg byte) (byte, error) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	_, err := bus.rc.Write([]byte{reg})
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 1)
	_, err = bus.rc.Read(buf)

	if err != nil {
		return 0, err
	}

	return buf[0], nil
}

func (bus *i2c) Write(reg byte, val byte) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	buf := []byte{reg, val}
	_, err := bus.rc.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (bus *i2c) ReadLen(reg byte, val []byte) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	_, err := bus.rc.Write([]byte{reg})
	if err != nil {
		return err
	}

	_, err = bus.rc.Read(val)
	if err != nil {
		return err
	}

	return nil
}

func (bus *i2c) WriteLen(reg byte, val []byte) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	_, err := bus.rc.Write(append([]byte{reg}, val...))
	if err != nil {
		return err
	}

	return nil
}

func (bus *i2c) Close() error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	return bus.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}

	return nil
}
