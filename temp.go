package raspberryheat

import (
	"os"
	"strconv"
	"strings"
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
		sensors = append(sensors, Sensor(W1Path+f.Name())+"/w1_slave")
	}
	return sensors
}

func (s Sensor) Read() float64 {
	var Buf [100]byte
	buf := Buf[:]
	f, err := os.Open(string(s))
	check(err)
	defer f.Close()

	n, _ := f.Read(buf)
	out := string(buf)
	start := strings.Index(out, "t=") + 2
	temp := out[start:n-1]
	t, err := strconv.ParseFloat(temp, 64)
	check(err)
	return t/1000
}
