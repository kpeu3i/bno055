package bno055

type CalibrationOffsets []byte

type CalibrationStatus struct {
	System        byte
	Gyroscope     byte
	Accelerometer byte
	Magnetometer  byte
}

func newCalibrationStatus(status byte) *CalibrationStatus {
	calibration := &CalibrationStatus{
		System:        (status >> 6) & 0x03,
		Gyroscope:     (status >> 4) & 0x03,
		Accelerometer: (status >> 2) & 0x03,
		Magnetometer:  status & 0x03,
	}

	return calibration
}

func (c *CalibrationStatus) IsCalibrated() bool {
	return c.System == 3 && c.Gyroscope == 3 && c.Accelerometer == 3 && c.Magnetometer == 3
}
