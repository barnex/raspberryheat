package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const W1Path = "/sys/bus/w1/devices/"

type Sensor string

func LsSensors() []*Sensor {
	devices, err := os.Open(W1Path)
	check(err)
	defer devices.Close()

	fi, err := devices.Readdir(-1)
	check(err)

	var sensors []*Sensor
	for _, f := range fi {
		if f.Name() == "w1_bus_master1" {
			continue
		}
		s := Sensor(W1Path + f.Name() + "/w1_slave")
		sensors = append(sensors, &s)
	}
	return sensors
}

func (s *Sensor) Read() (float64, error) {
	var Buf [100]byte
	buf := Buf[:]
	f, err := os.Open(string(*s))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err2 := f.Read(buf)
	if err2 != nil {
		return 0, err2
	}
	out := string(buf)
	crcfail := strings.Index(out, "NO")
	if crcfail != -1 {
		return 0, errors.New("CRC fail")
	}
	start := strings.Index(out, "t=") + 2
	temp := out[start : n-1]
	t, err3 := strconv.ParseFloat(temp, 64)
	if err3 != nil {
		return 0, err3
	}

	return t / 1000, nil
}
