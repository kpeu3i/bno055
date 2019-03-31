package bno055

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

const (
	i2cSlave = 0x0703
)

type i2c struct {
	retryCount   int
	retryTimeout time.Duration
	mu           sync.Mutex
	rc           *os.File
}

func newI2C(addr uint8, bus int, retryCount int, retryTimeout time.Duration) (*i2c, error) {
	file, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	err = ioctl(file.Fd(), i2cSlave, uintptr(addr))
	if err != nil {
		return nil, err
	}

	i2c := &i2c{
		retryCount:   retryCount,
		retryTimeout: retryTimeout,
		rc:           file,
	}

	return i2c, nil
}

func (b *i2c) Read(reg byte) (byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	buf := make([]byte, 1)

	err := retry(func() error {
		_, err := b.rc.Write([]byte{reg})
		if err != nil {
			return err
		}

		_, err = b.rc.Read(buf)

		if err != nil {
			return err
		}

		return nil
	}, b.retryCount, b.retryTimeout)

	return buf[0], err
}

func (b *i2c) Write(reg byte, val byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	buf := []byte{reg, val}

	err := retry(func() error {
		_, err := b.rc.Write(buf)
		if err != nil {
			return err
		}

		return nil
	}, b.retryCount, b.retryTimeout)

	return err
}

func (b *i2c) ReadLen(reg byte, val []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := retry(func() error {
		_, err := b.rc.Write([]byte{reg})
		if err != nil {
			return err
		}

		_, err = b.rc.Read(val)
		if err != nil {
			return err
		}

		return nil
	}, b.retryCount, b.retryTimeout)

	return err
}

func (b *i2c) WriteLen(reg byte, val []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := retry(func() error {
		_, err := b.rc.Write(append([]byte{reg}, val...))
		if err != nil {
			return err
		}

		return nil
	}, b.retryCount, b.retryTimeout)

	return err
}

func (b *i2c) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.rc.Close()
}

func retry(fn func() error, retryCount int, timeout time.Duration) error {
	var err error

	for i := 1; i <= retryCount+1; i++ {
		err = fn()
		if err == nil {
			break
		}

		time.Sleep(timeout)
	}

	return err
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}

	return nil
}
