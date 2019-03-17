# BNO055

A userspace I²C driver for the Bosch [BNO055](https://www.bosch-sensortec.com/bst/products/all_products/bno055) 9-axis Absolute Orientation Sensor.

## Install

Use go get to install the latest version of the library:

    go get github.com/kpeu3i/bno055@v1.0.0

Next, include bno055 in your application:

```go
import "github.com/kpeu3i/bno055"
```

## Usage

First, connect the sensor via I²C interface.

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kpeu3i/bno055"
)

func main() {
	sensor, err := bno055.NewSensor(0x28, 1)
	if err != nil {
		panic(err)
	}

	err = sensor.UseExternalCrystal(true)
	if err != nil {
		panic(err)
	}

	status, err := sensor.Status()
	if err != nil {
		panic(err)
	}

	fmt.Printf("*** Status: system=%v, system_error=%v, self_test=%v\n", status.System, status.SystemError, status.SelfTest)

	revision, err := sensor.Revision()
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"*** Revision: software=%v, bootloader=%v, accelerometer=%v, gyroscope=%v, magnetometer=%v\n",
		revision.Software,
		revision.Bootloader,
		revision.Accelerometer,
		revision.Gyroscope,
		revision.Magnetometer,
	)

	axisConfig, err := sensor.AxisConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"*** Axis: x=%v, y=%v, z=%v, sign_x=%v, sign_y=%v, sign_z=%v\n",
		axisConfig.X,
		axisConfig.Y,
		axisConfig.Z,
		axisConfig.SignX,
		axisConfig.SignY,
		axisConfig.SignZ,
	)

	temperature, err := sensor.Temperature()
	if err != nil {
		panic(err)
	}

	fmt.Printf("*** Temperature: t=%v\n", temperature)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			err := sensor.Close()
			if err != nil {
				panic(err)
			}
		default:
			vector, err := sensor.Euler()
			if err != nil {
				panic(err)
			}

			fmt.Printf("\r*** Euler angles: x=%5.3f, y=%5.3f, z=%5.3f", vector.X, vector.Y, vector.Z)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Output:
	// *** Status: system=133, system_error=0, self_test=15
	// *** Revision: software=785, bootloader=21, accelerometer=251, gyroscope=15, magnetometer=50
	// *** Axis: x=0, y=1, z=2, sign_x=0, sign_y=0, sign_z=0
	// *** Temperature: t=27
	// *** Euler angles: x=2.312, y=2.000, z=91.688
}
```

## Troubleshooting

#### How to enable I²C bus on RPi device?

If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS level.
Go to "Interfacing Options" menu, to active I²C bus.
Probably you will need to reboot to load I²C kernel module.
Finally you should have device like `/dev/i2c-1` present in the system.

#### How to find I²C bus allocation and device address?

Use `i2cdetect` utility in format "i2cdetect -y X", where X may vary from 0 to 5 or more,
to discover address occupied by peripheral device. To install utility you should run
`apt install i2c-tools` on debian-kind system. `i2cdetect -y 1` sample output:

```
         0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
    00:          -- -- -- -- -- -- -- -- -- -- -- -- --
    10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    20: -- -- -- -- -- -- -- -- 28 -- -- -- -- -- -- --
    30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    70: -- -- -- -- -- -- -- --

```

#### Workaround for I²C clock stretching

It seems all versions of Raspberry Pi have an I²C bus [hardware problem](http://www.advamation.com/knowhow/raspberrypi/rpi-i2c-bug.html) preventing them from working correctly with Bosch BNO055.
The problem has been variously diagnosed as being due to the Pi’s inability to handle clock stretching in arbitrary parts of the I²C transaction and the BNO055 chip’s exquisite sensitivity to I²C bus levels.

Solutions:

1. [Configuring software I²C driver](https://github.com/fivdi/i2c-bus/blob/master/doc/raspberry-pi-software-i2c.md)

    Raspbian has a software I²C driver that can be enabled by adding the following line to `/boot/config.txt`:
    ```
    dtoverlay=i2c-gpio,bus=3
    ```
    This will create an I²C bus called `/dev/i2c-3`. SDA will be on GPIO23 and SCL will be on GPIO24 which are pins 16 and 18 on the GPIO header respectively.

2. [Slowing the I²C bus transactions](https://softsolder.com/2018/08/22/raspberry-pi-3-i2c-vs-bosch-bno055-absolute-orientation-sensor/)

    The solution require slowing the I²C bus transactions to 25 kb/s, by inserting a line in the `/boot/config.txt` file:
    ```
    dtparam=i2c_arm_baudrate=25000
    ```

## TODO

* Docs
* Tests
