package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	W1Path    = "/sys/bus/w1/devices/"
	MaxErrors = 10 // maximum number of successive read errors before sensor is considered faulty
)

type Sensor struct {
	file    string  // sensor device file
	temp    float64 // room temperature in °C
	led     *GPIO   // room LED to blink
	sumTemp float64 // to track average temperature
	sumN    int     // to track average temperature
	errRun  int     // number of successive errors, too many -> temp=NaN
	err     error   // last error, if any
	sync.Mutex
}

func NewSensor(file string, led *GPIO) *Sensor {
	return &Sensor{file: W1Path + file + "/w1_slave", led: led, temp: math.NaN()}
}

func (s *Sensor) AvgTemp() float64 {
	s.Lock()
	defer s.Unlock()
	avg := s.sumTemp / float64(s.sumN)
	s.sumTemp = 0
	s.sumN = 0
	return avg
}

func (s *Sensor) Error() string {
	s.Lock()
	defer s.Unlock()
	if s.err != nil {
		return s.err.Error()
	}
	return ""
}

func (s *Sensor) Update() {
	t, err := s.Read()
	s.Lock()

	// report temp read error
	if err != nil {
		s.err = err
		s.errRun++
		if s.errRun > MaxErrors {
			s.temp = math.NaN()
		}
		s.Unlock()
		blinkErr(s.led)
		return
	}

	// store temp
	s.temp = t
	s.errRun, s.err = 0, nil
	// track average
	s.sumN++
	s.sumTemp += t

	s.Unlock()
	blinkOK(s.led) // report success
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
		return 0, fmt.Errorf("CRC fail")
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
		return 0, fmt.Errorf("measurment returned too fast")
	}

	// disconnected 3V3 can just give wrong temp
	// so do sanity check
	if t > 40*1000 || t < -10*1000 {
		return 0, fmt.Errorf("measurment too extreme: %v", t)
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
