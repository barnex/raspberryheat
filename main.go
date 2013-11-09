package main

import (
	"time"
)

var (
	sensor = []*Sensor{
		NewSensor("28-0000050ad012", "living", LED1),
		NewSensor("28-0000050b07f7", "kindjes", LED2)}
)

func main() {
	go StartHTTP()

	for {
		for _, s := range sensor {
			s.Update()
		}
	}
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
