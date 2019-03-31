// Package bno055 allows interfacing with the BNO055 absolute orientation sensor.
package bno055

import (
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

type Status struct {
	// System Status (see section 4.3.58)
	// ---------------------------------
	// 0 = Idle
	// 1 = System Error
	// 2 = Initializing Peripherals
	// 3 = System Initialization
	// 4 = Executing Self-Test
	// 5 = Sensor fusion algorithm running
	// 6 = System running without fusion algorithms
	System uint8

	// System Error (see section 4.3.59)
	//---------------------------------
	// 0 = No error
	// 1 = Peripheral initialization error
	// 2 = System initialization error
	// 3 = Self test result failed
	// 4 = Register map value out of range
	// 5 = Register map address out of range
	// 6 = Register map write error
	// 7 = BNO low power mode not available for selected operation ion mode
	// 8 = Accelerometer power mode not available
	// 9 = Fusion algorithm configuration error
	// A = Sensor configuration error
	SystemError uint8

	// Self Test Results
	// --------------------------------
	// 1 = test passed, 0 = test failed
	//
	// Bit 0 = Accelerometer self test
	// Bit 1 = Magnetometer self test
	// Bit 2 = Gyroscope self test
	// Bit 3 = MCU self test
	//
	// 0x0F = all good!
	SelfTest uint8
}

type Revision struct {
	Software      uint16
	Bootloader    uint8
	Gyroscope     uint8
	Accelerometer uint8
	Magnetometer  uint8
}

type Vector struct {
	X float32
	Y float32
	Z float32
}

type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

type Option func(sensor *Sensor)

var defaultCalibrationOffsets CalibrationOffsets = []byte{
	239, 255, 184, 255, 10, 0, 196, 0, 193, 0,
	85, 255, 128, 0, 0, 0, 1, 0, 232, 3, 0, 0,
}

type Sensor struct {
	retryCount   int
	retryTimeout time.Duration
	bus          *i2c
	mu           sync.Mutex
	opMode       byte
}

func (s *Sensor) Status() (*Status, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.bus.Write(bno055PageID, 0x0)
	if err != nil {
		return nil, err
	}

	prevMode := s.opMode

	err = s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return nil, err
	}

	sysTrigger, err := s.bus.Read(bno055SysTrigger)
	err = s.bus.Write(bno055SysTrigger, sysTrigger|0x1)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	err = s.setOperationMode(prevMode)
	if err != nil {
		return nil, err
	}

	system, err := s.bus.Read(bno055SysStat)
	if err != nil {
		return nil, err
	}

	selfTest, err := s.bus.Read(bno055SelfTestResult)
	if err != nil {
		return nil, err
	}

	systemError, err := s.bus.Read(bno055SysErr)
	if err != nil {
		return nil, err
	}

	status := &Status{
		System:      system,
		SystemError: systemError,
		SelfTest:    selfTest,
	}

	return status, nil
}

func (s *Sensor) Revision() (*Revision, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	accelerometer, err := s.bus.Read(bno055AccelRevID)
	if err != nil {
		return nil, err
	}

	magnetometer, err := s.bus.Read(bno055MagRevID)
	if err != nil {
		return nil, err
	}

	gyroscope, err := s.bus.Read(bno055GyroRevID)
	if err != nil {
		return nil, err
	}

	bootloader, err := s.bus.Read(bno055BLRevID)
	if err != nil {
		return nil, err
	}

	swLSB, err := s.bus.Read(bno055SWRevIDLsb)
	if err != nil {
		return nil, err
	}

	swMSB, err := s.bus.Read(bno055SWRevIDMsb)
	if err != nil {
		return nil, err
	}

	software := (uint16(swMSB) << 8) | uint16(swLSB)

	revision := &Revision{
		Accelerometer: accelerometer,
		Magnetometer:  magnetometer,
		Gyroscope:     gyroscope,
		Bootloader:    bootloader,
		Software:      software,
	}

	return revision, err
}

func (s *Sensor) UseExternalCrystal(b bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prevMode := s.opMode

	err := s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return err
	}

	err = s.bus.Write(bno055PageID, 0)
	if err != nil {
		return err
	}

	if b {
		err = s.bus.Write(bno055SysTrigger, 0x80)
		if err != nil {
			return err
		}
	} else {
		err = s.bus.Write(bno055SysTrigger, 0x00)
		if err != nil {
			return err
		}
	}

	err = s.setOperationMode(prevMode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sensor) Calibration() (CalibrationOffsets, *CalibrationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, err := s.bus.Read(bno055CalibStat)
	if err != nil {
		return nil, nil, err
	}

	prevMode := s.opMode

	err = s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return nil, nil, err
	}

	offsets := make([]byte, 22)
	err = s.bus.ReadLen(bno055AccelOffsetXLsb, offsets)
	if err != nil {
		return nil, nil, err
	}

	err = s.setOperationMode(prevMode)
	if err != nil {
		return nil, nil, err
	}

	calibrationOffsets := CalibrationOffsets(offsets)
	calibrationStatus := newCalibrationStatus(status)

	return calibrationOffsets, calibrationStatus, nil
}

func (s *Sensor) Calibrate(offsets CalibrationOffsets) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prevMode := s.opMode

	err := s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return err
	}

	err = s.bus.WriteLen(bno055AccelOffsetXLsb, offsets)
	if err != nil {
		return err
	}

	err = s.setOperationMode(prevMode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sensor) AxisConfig() (*AxisConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mapConfig, err := s.bus.Read(bno055AxisMapConfig)
	if err != nil {
		return nil, err
	}

	signConfig, err := s.bus.Read(bno055AxisMapSign)
	if err != nil {
		return nil, err
	}

	axisConfig := newAxisConfig(mapConfig, signConfig)

	return axisConfig, nil
}

// Note that by default the axis orientation of the BNO chip looks like the
// following (taken from section 3.4, page 24 of the datasheet).
// Notice the dot in the corner that corresponds to the dot on the BNO chip:
//
//                   | Z axis
//                   |
//                   |   / X axis
//               ____|__/____
//  Y axis     / *   | /    /|
//  _________ /______|/    //
//           /___________ //
//          |____________|/
//
func (s *Sensor) RemapAxis(config *AxisConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prevMode := s.opMode

	err := s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return err
	}

	err = s.bus.Write(bno055AxisMapConfig, config.Mappings())
	if err != nil {
		return err
	}

	err = s.bus.Write(bno055AxisMapSign, config.Signs())
	if err != nil {
		return err
	}

	err = s.setOperationMode(prevMode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sensor) Temperature() (int8, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	temperature, err := s.bus.Read(bno055Temp)
	if err != nil {
		return 0, err
	}

	return int8(temperature), nil
}

func (s *Sensor) Magnetometer() (*Vector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	x, y, z, err := s.readVector(bno055MagDataXLsb)
	if err != nil {
		return nil, err
	}

	// 1uT = 16 LSB
	vector := &Vector{
		X: float32(x) / 16,
		Y: float32(y) / 16,
		Z: float32(z) / 16,
	}

	return vector, nil
}

func (s *Sensor) Gyroscope() (*Vector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	x, y, z, err := s.readVector(bno055GyroDataXLsb)
	if err != nil {
		return nil, err
	}

	// 1 degree = 16 LSB
	vector := &Vector{
		X: float32(x) / 16,
		Y: float32(y) / 16,
		Z: float32(z) / 16,
	}

	return vector, nil
}

func (s *Sensor) Euler() (*Vector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	x, y, z, err := s.readVector(bno055EulerHLsb)
	if err != nil {
		return nil, err
	}

	// 1rps = 16 LSB
	vector := &Vector{
		X: float32(x) / 16,
		Y: float32(y) / 16,
		Z: float32(z) / 16,
	}

	return vector, nil
}

func (s *Sensor) Accelerometer() (*Vector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	x, y, z, err := s.readVector(bno055AccelDataXLsb)
	if err != nil {
		return nil, err
	}

	// 1m/s^2 = 100 LSB
	vector := &Vector{
		X: float32(x) / 100,
		Y: float32(y) / 100,
		Z: float32(z) / 100,
	}

	return vector, nil
}

func (s *Sensor) LinearAccelerometer() (*Vector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	x, y, z, err := s.readVector(bno055LinearAccelDataXLsb)
	if err != nil {
		return nil, err
	}

	// 1m/s^2 = 100 LSB
	vector := &Vector{
		X: float32(x) / 100,
		Y: float32(y) / 100,
		Z: float32(z) / 100,
	}

	return vector, nil
}

func (s *Sensor) Gravity() (*Vector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	x, y, z, err := s.readVector(bno055GravityDataXLsb)
	if err != nil {
		return nil, err
	}

	// 1m/s^2 = 100 LSB
	vector := &Vector{
		X: float32(x) / 100,
		Y: float32(y) / 100,
		Z: float32(z) / 100,
	}

	return vector, nil
}

func (s *Sensor) Quaternion() (*Quaternion, error) {
	w, x, y, z, err := s.readQuaternion(bno055QuaternionDataWLsb)
	if err != nil {
		return nil, err
	}

	scale := float32(1 / (1 << 14))

	quaternion := &Quaternion{
		W: scale * float32(w),
		X: scale * float32(x),
		Y: scale * float32(y),
		Z: scale * float32(z),
	}

	return quaternion, nil
}

func (s *Sensor) Sleep() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prevMode := s.opMode

	err := s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return err
	}

	err = s.bus.Write(bno055PwrMode, bno055PowerModeSuspend)
	if err != nil {
		return err
	}

	err = s.setOperationMode(prevMode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sensor) Wakeup() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prevMode := s.opMode

	err := s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return err
	}

	err = s.bus.Write(bno055PwrMode, bno055PowerModeNormal)
	if err != nil {
		return err
	}

	err = s.setOperationMode(prevMode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sensor) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.bus.Close()
}

func (s *Sensor) setOperationMode(mode byte) error {
	err := s.bus.Write(bno055OprMode, mode)
	if err != nil {
		return err
	}

	s.opMode = mode

	return nil
}

func (s *Sensor) readVector(addr byte) (x, y, z int16, err error) {
	buf := make([]byte, 6)
	err = s.bus.ReadLen(addr, buf)
	if err != nil {
		return
	}

	x = int16(binary.LittleEndian.Uint16(buf[0:]))
	y = int16(binary.LittleEndian.Uint16(buf[2:]))
	z = int16(binary.LittleEndian.Uint16(buf[4:]))

	return
}

func (s *Sensor) readQuaternion(addr byte) (w, x, y, z int16, err error) {
	buf := make([]byte, 8)
	err = s.bus.ReadLen(addr, buf)
	if err != nil {
		return
	}

	w = int16(binary.LittleEndian.Uint16(buf[0:]))
	x = int16(binary.LittleEndian.Uint16(buf[2:]))
	y = int16(binary.LittleEndian.Uint16(buf[4:]))
	z = int16(binary.LittleEndian.Uint16(buf[6:]))

	return
}

func (s *Sensor) checkExists() error {
	for i := 0; i < 10; i++ {
		id, err := s.bus.Read(bno055ChipID)
		if err != nil {
			return err
		}

		if id == bno055Id {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return errors.New("sensor not found")
}

func (s *Sensor) init() error {
	err := s.checkExists()
	if err != nil {
		return err
	}

	err = s.setOperationMode(bno055OperationModeConfig)
	if err != nil {
		return err
	}

	// Reset the device using the reset command
	err = s.bus.Write(bno055SysTrigger, 0x20)
	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond)

	err = s.checkExists()
	if err != nil {
		return err
	}

	// Set to normal power mode
	err = s.bus.Write(bno055PwrMode, bno055PowerModeNormal)
	if err != nil {
		return err
	}

	err = s.bus.Write(bno055PageID, 0x0)
	if err != nil {
		return err
	}

	// Default to internal oscillator
	err = s.bus.Write(bno055SysTrigger, 0x00)
	if err != nil {
		return err
	}

	// Set temperature source to gyroscope, as it seems to be more accurate
	err = s.bus.Write(bno055TempSource, 0x01)
	if err != nil {
		return err
	}

	// Set the unit selection bits
	err = s.bus.Write(bno055UnitSel, 0x0)
	if err != nil {
		return err
	}

	err = s.setOperationMode(bno055OperationModeNdof)
	if err != nil {
		return err
	}

	return nil
}

func WithRetry(retryCount int, retryTimeout time.Duration) Option {
	return func(sensor *Sensor) {
		sensor.retryCount = retryCount
		sensor.retryTimeout = retryTimeout
	}
}

func NewSensor(addr uint8, bus int, options ...Option) (*Sensor, error) {
	sensor := &Sensor{opMode: bno055OperationModeNdof}
	for _, option := range options {
		option(sensor)
	}

	i2c, err := newI2C(addr, bus, sensor.retryCount, sensor.retryTimeout)
	if err != nil {
		return nil, err
	}

	sensor.bus = i2c

	err = sensor.init()
	if err != nil {
		return nil, err
	}

	err = sensor.Calibrate(defaultCalibrationOffsets)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}
