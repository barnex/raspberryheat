package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const W1Path = "/sys/bus/w1/devices/"

type Sensor struct {
	file        string
	description string
	temp        float64
	led         *GPIO
	sync.Mutex
	err error
}

func NewSensor(file, name string, led *GPIO) *Sensor {
	return &Sensor{file: W1Path + file + "/w1_slave", description: name, led: led}
}

func (s *Sensor) Label(prefix string) string {
	return prefix + "_" + s.description
}

func (s *Sensor) Error() string {
	s.Lock()
	defer s.Unlock()
	if s.err != nil {
		return s.err.Error()
	}
	return ""
}

func (s *Sensor) Description() string {
	return s.description
}

func (s *Sensor) Update() {
	t, err := s.Read()

	s.Lock()

	// report temp read error
	if err != nil {
		s.err = err
		s.Unlock()
		blinkErr(s.led)
		return
	}

	// store temp
	s.temp = t
	s.err = nil
	s.Unlock()

	// report success
	blinkOK(s.led)
}

func (s *Sensor) Temp() float64 {
	s.Lock()
	defer s.Unlock()
	return s.temp
}

func (s *Sensor) Read() (float64, error) {
	var Buf [100]byte
	buf := Buf[:]

	startt := time.Now()

	// try to open device file
	f, err := os.Open(string(s.file))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// try to read device file
	n, err2 := f.Read(buf)
	if err2 != nil {
		return 0, err2
	}

	// look for failed CRC
	out := string(buf)
	crcfail := strings.Index(out, "NO")
	if crcfail != -1 {
		return 0, fmt.Errorf("%v CRC fail", s.description)
	}

	// parse temperature (milligrades)
	start := strings.Index(out, "t=") + 2
	temp := out[start : n-1]
	t, err3 := strconv.ParseFloat(temp, 64)
	if err3 != nil {
		return 0, err3
	}

	// disconnected 3V3 can give wrong temp, correct CRC but returns quickly
	if time.Since(startt) < 500*time.Millisecond {
		return 0, fmt.Errorf("%v returned too fast", s.description)
	}

	// disconnected 3V3 can just give wrong temp
	// so do sanity check
	if t > 40*1000 || t < -10*1000 {
		return 0, fmt.Errorf("%v temperature too extreme: %v", s.description, t)
	}

	return t / 1000, nil
}

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
		s := &Sensor{file: W1Path + f.Name() + "/w1_slave"}
		sensors = append(sensors, s)
	}
	return sensors
}
