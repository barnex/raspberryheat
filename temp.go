package raspberryheat

import (
	"os"
)

const W1Path = "/sys/bus/w1/devices/"

type Sensor string

func LsSensors() []Sensor {
	devices, err := os.Open(W1Path)
	check(err)
	defer devices.Close()

	fi, err := devices.Readdir(-1)
	check(err)

	var sensors []Sensor
	for _, f := range fi {
		if f.Name() == "w1_bus_master1" {
			continue
		}
		sensors = append(sensors, Sensor(W1Path+f.Name()))
	}
	return sensors
}
