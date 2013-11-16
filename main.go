package main

import (
	"time"
)

const logPeriod = 2 * time.Minute

var HEATER = RELAY2

var rooms = []*Room{
	NewRoom("living", "28-0000050ad012", LED1),
	NewRoom("kindjes", "28-0000050b07f7", LED2)}

var Burn bool

func main() {
	go StartHTTP()

	lastLog := time.Now()
	for {

		for _, r := range rooms {
			r.Update()
		}

		Burn = false
		for _, r := range rooms {
			if r.Burn{
				Burn = true
				break
			}
		}
		HEATER.Set(Burn)

		if time.Since(lastLog) > logPeriod {
			doLog()
			lastLog = time.Now()
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
