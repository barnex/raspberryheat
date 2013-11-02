package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	statusLED  = LED1
	measureLED = LED2
	httpLED    = LED3
	errorLED   = LED4
)

var (
	temp        []float64
	measureLock sync.Mutex
)

func getTemp() []float64 {
	measureLock.Lock()
	defer measureLock.Unlock()
	return temp
}

func main() {

	sensors := LsSensors()
	log.Println(sensors)
	temp = make([]float64, len(sensors))

	go StartHTTP()

	for {
		blink(statusLED)

		for i, s := range sensors {
			t, err := s.Read()
			if err != nil {
				blink(errorLED)
				log.Println(err)
				continue
			}
			measureLock.Lock()
			temp[i] = t
			measureLock.Unlock()
			fmt.Print(t, " ")
			blink(measureLED)
		}
		fmt.Println()
	}
}

func blink(led *GPIO) {
	led.Set(true)
	time.Sleep(10 * time.Millisecond)
	led.Set(false)
}
