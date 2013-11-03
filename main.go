package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	led         = []*GPIO{LED1, LED2}
	temp        = []float64{0, 0}
	sensors     = []*Sensor{NewSensor("28-0000050ad012"), NewSensor("28-0000050b07f7")}
	roomName    = []string{"living", "kindjes"}
	measureLock sync.Mutex
)

func getTemp() []float64 {
	measureLock.Lock()
	defer measureLock.Unlock()
	return temp
}

func main() {

	go StartHTTP()

	for {
		measure()
	}
}

var (
	errTooFast  = errors.New("measurement returned too fast")
	errTooLarge = errors.New("temperature too extreme")
)

func measure() {
	for i, s := range sensors {

		// disconnected 3V3 can give wrong temp, correct CRC but returns quickly
		start := time.Now()
		t, err := s.Read()
		if err == nil && time.Since(start) < 500*time.Millisecond {
			err = errTooFast
		}

		// disconnected 3V3 can just give wrong temp
		if err == nil && t > 40 || t < -10 {
			err = errTooLarge
		}

		// report temp read error
		if err != nil {
			blinkErr(led[i])
			log.Println(err)
			continue
		}

		// store temp
		measureLock.Lock()
		temp[i] = t
		measureLock.Unlock()

		// report success
		fmt.Print(t, " ")
		blinkOK(led[i])
	}
	fmt.Println()
}

// short blink indicates successful measurement
func blinkOK(led *GPIO) {
	led.Set(true)
	time.Sleep(10 * time.Millisecond)
	led.Set(false)
}

// briefly off indicates measure error
func blinkErr(led *GPIO) {
	led.Set(false)
	time.Sleep(50 * time.Millisecond)
	led.Set(true)
}
