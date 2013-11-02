//+build ignore

package main

import (
	. "."
	"log"
	"sync"
	"fmt"
	"time"
)

var (
	statusLED  = LED1
	measureLED = LED2
	errorLED   = LED3
)

var (
	temp        []float64
	measureLock sync.Mutex
)

func main() {

	sensors := LsSensors()
	log.Println(sensors)
	temp = make([]float64, len(sensors))

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
