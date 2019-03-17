package bno055

type AxisConfig struct {
	X     byte
	Y     byte
	Z     byte
	SignX byte
	SignY byte
	SignZ byte
}

func newAxisConfig(mapConfig, signConfig byte) *AxisConfig {
	axisConfig := &AxisConfig{
		X:     mapConfig & 0x03,
		Y:     (mapConfig >> 2) & 0x03,
		Z:     (mapConfig >> 4) & 0x03,
		SignX: (signConfig >> 2) & 0x01,
		SignY: (signConfig >> 1) & 0x01,
		SignZ: signConfig & 0x01,
	}

	return axisConfig
}

func (c *AxisConfig) Mappings() byte {
	var mappings byte

	mappings |= (c.Z & 0x03) << 4
	mappings |= (c.Y & 0x03) << 2
	mappings |= c.X & 0x03

	return mappings
}

func (c *AxisConfig) Signs() byte {
	var signs byte

	signs |= (c.SignX & 0x01) << 2
	signs |= (c.SignY & 0x01) << 1
	signs |= c.SignZ & 0x01

	return signs
}
