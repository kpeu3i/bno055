package bno055

const (
	// Constant device identifier
	bno055Id = 0xA0

	// Page id register definition
	bno055PageID = 0x07

	// PAGE0 register definition start
	bno055ChipID     = 0x00
	bno055AccelRevID = 0x01
	bno055MagRevID   = 0x02
	bno055GyroRevID  = 0x03
	bno055SWRevIDLsb = 0x04
	bno055SWRevIDMsb = 0x05
	bno055BLRevID    = 0x06

	// Accelerometer data register
	bno055AccelDataXLsb = 0x08
	bno055AccelDataXMsb = 0x09
	bno055AccelDataYLsb = 0x0A
	bno055AccelDataYMsb = 0x0B
	bno055AccelDataZLsb = 0x0C
	bno055AccelDataZMsb = 0x0D

	// Magnetometer data register
	bno055MagDataXLsb = 0x0E
	bno055MagDataXMsb = 0x0F
	bno055MagDataYLsb = 0x10
	bno055MagDataYMsb = 0x11
	bno055MagDataZLsb = 0x12
	bno055MagDataZMsb = 0x13

	// Gyroscope data registers
	bno055GyroDataXLsb = 0x14
	bno055GyroDataXMsb = 0x15
	bno055GyroDataYLsb = 0x16
	bno055GyroDataYMsb = 0x17
	bno055GyroDataZLsb = 0x18
	bno055GyroDataZMsb = 0x19

	// Euler data registers
	bno055EulerHLsb = 0x1A
	bno055EulerHMsb = 0x1B
	bno055EulerRLsb = 0x1C
	bno055EulerRMsb = 0x1D
	bno055EulerPLsb = 0x1E
	bno055EulerPMsb = 0x1F

	// Quaternion data registers
	bno055QuaternionDataWLsb = 0x20
	bno055QuaternionDataWMsb = 0x21
	bno055QuaternionDataXLsb = 0x22
	bno055QuaternionDataXMsb = 0x23
	bno055QuaternionDataYLsb = 0x24
	bno055QuaternionDataYMsb = 0x25
	bno055QuaternionDataZLsb = 0x26
	bno055QuaternionDataZMsb = 0x27

	// Linear acceleration data registers
	bno055LinearAccelDataXLsb = 0x28
	bno055LinearAccelDataXMsb = 0x29
	bno055LinearAccelDataYLsb = 0x2A
	bno055LinearAccelDataYMsb = 0x2B
	bno055LinearAccelDataZLsb = 0x2C
	bno055LinearAccelDataZMsb = 0x2D

	// Gravity data registers
	bno055GravityDataXLsb = 0x2E
	bno055GravityDataXMsb = 0x2F
	bno055GravityDataYLsb = 0x30
	bno055GravityDataYMsb = 0x31
	bno055GravityDataZLsb = 0x32
	bno055GravityDataZMsb = 0x33

	// Temperature data register
	bno055Temp = 0x34

	// Status registers
	bno055CalibStat      = 0x35
	bno055SelfTestResult = 0x36
	bno055IntrStat       = 0x37

	bno055SysClkStat = 0x38
	bno055SysStat    = 0x39
	bno055SysErr     = 0x3A

	// Unit selection register
	bno055UnitSel    = 0x3B
	bno055AataSelect = 0x3C

	// Mode registers
	bno055OprMode = 0x3D
	bno055PwrMode = 0x3E

	bno055SysTrigger = 0x3F
	bno055TempSource = 0x40

	// AxisConfig remap registers
	bno055AxisMapConfig = 0x41
	bno055AxisMapSign   = 0x42

	// SIC registers
	bno055SicMatrix0Lsb = 0x43
	bno055SicMatrix0Msb = 0x44
	bno055SicMatrix1Lsb = 0x45
	bno055SicMatrix1Msb = 0x46
	bno055SicMatrix2Lsb = 0x47
	bno055SicMatrix2Msb = 0x48
	bno055SicMatrix3Lsb = 0x49
	bno055SicMatrix3Msb = 0x4A
	bno055SicMatrix4Lsb = 0x4B
	bno055SicMatrix4Msb = 0x4C
	bno055SicMatrix5Lsb = 0x4D
	bno055SicMatrix5Msb = 0x4E
	bno055SicMatrix6Lsb = 0x4F
	bno055SicMatrix6Msb = 0x50
	bno055SicMatrix7Lsb = 0x51
	bno055SicMatrix7Msb = 0x52
	bno055SicMatrix8Lsb = 0x53
	bno055SicMatrix8Msb = 0x54

	// Accelerometer offset registers
	bno055AccelOffsetXLsb = 0x55
	bno055AccelOffsetXMsb = 0x56
	bno055AccelOffsetYLsb = 0x57
	bno055AccelOffsetYMsb = 0x58
	bno055AccelOffsetZLsb = 0x59
	bno055AccelOffsetZMsb = 0x5A

	// Magnetometer offset registers
	bno055MagOffsetXLsb = 0x5B
	bno055MagOffsetXMsb = 0x5C
	bno055MagOffsetYLsb = 0x5D
	bno055MagOffsetYMsb = 0x5E
	bno055MagOffsetZLsb = 0x5F
	bno055MagOffsetZMsb = 0x60

	// Gyroscope offset registers
	bno055GyroOffsetXLsb = 0x61
	bno055GyroOffsetXMsb = 0x62
	bno055GyroOffsetYLsb = 0x63
	bno055GyroOffsetYMsb = 0x64
	bno055GyroOffsetZLsb = 0x65
	bno055GyroOffsetZMsb = 0x66

	// Radius registers
	bno055AccelRadiusLsb = 0x67
	bno055AccelRadiusMsb = 0x68
	bno055MagRadiusLsb   = 0x69
	bno055MagRadiusMsb   = 0x6A

	bno055PowerModeNormal   = 0x00
	bno055PowerModeLowpower = 0x01
	bno055PowerModeSuspend  = 0x02

	// Operation mode settings
	bno055OperationModeConfig     = 0x00
	bno055OperationModeAcconly    = 0x01
	bno055OperationModeMagonly    = 0x02
	bno055OperationModeGyronly    = 0x03
	bno055OperationModeAccmag     = 0x04
	bno055OperationModeAccgyro    = 0x05
	bno055OperationModeMaggyro    = 0x06
	bno055OperationModeAmg        = 0x07
	bno055OperationModeImuplus    = 0x08
	bno055OperationModeCompass    = 0x09
	bno055OperationModeM4g        = 0x0A
	bno055OperationModeNdofFmcOff = 0x0B
	bno055OperationModeNdof       = 0x0C

	bno055RemapConfigP0 = 0x21
	bno055RemapConfigP1 = 0x24 // default
	bno055RemapConfigP2 = 0x24
	bno055RemapConfigP3 = 0x21
	bno055RemapConfigP4 = 0x24
	bno055RemapConfigP5 = 0x21
	bno055RemapConfigP6 = 0x21
	bno055RemapConfigP7 = 0x24

	bno055RemapSignP0 = 0x04
	bno055RemapSignP1 = 0x00 // default
	bno055RemapSignP2 = 0x06
	bno055RemapSignP3 = 0x02
	bno055RemapSignP4 = 0x03
	bno055RemapSignP5 = 0x01
	bno055RemapSignP6 = 0x07
	bno055RemapSignP7 = 0x05
)
