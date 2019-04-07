package i2c

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

type Bus struct {
	retryCount   int
	retryTimeout time.Duration
	mu           sync.Mutex
	rc           *os.File
}

func NewBus(addr uint8, bus int, retryCount int, retryTimeout time.Duration) (*Bus, error) {
	file, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	err = ioctl(file.Fd(), i2cSlave, uintptr(addr))
	if err != nil {
		return nil, err
	}

	i2cBus := &Bus{
		retryCount:   retryCount,
		retryTimeout: retryTimeout,
		rc:           file,
	}

	return i2cBus, nil
}

func (b *Bus) Read(reg byte) (byte, error) {
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

func (b *Bus) Write(reg byte, val byte) error {
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

func (b *Bus) ReadBuffer(reg byte, buff []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := retry(func() error {
		_, err := b.rc.Write([]byte{reg})
		if err != nil {
			return err
		}

		_, err = b.rc.Read(buff)
		if err != nil {
			return err
		}

		return nil
	}, b.retryCount, b.retryTimeout)

	return err
}

func (b *Bus) WriteBuffer(reg byte, buff []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := retry(func() error {
		_, err := b.rc.Write(append([]byte{reg}, buff...))
		if err != nil {
			return err
		}

		return nil
	}, b.retryCount, b.retryTimeout)

	return err
}

func (b *Bus) Close() error {
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
