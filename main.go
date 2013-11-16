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

		// Update all sensors
		for _, r := range rooms {
			r.sensor.Update()
		}

		// If at least one room is overheated, then we reset the other's
		// Schmidt triggers so they may stop heating uless absolutely required.
		var overheated bool
		for _, r := range rooms {
			if r.Overheat() {
				overheated = true
				break
			}
		}
		if overheated {
			for _, r := range rooms {
				// Reset Schmidt trigger unless right at the edge.
				if r.sensor.Temp() > r.SetTemp -r.Schmidt + 0.1 {
					r.Burn = false
				}
			}
		}

		Burn = false
		for _, r := range rooms {
			r.UpdateBurn()
			if r.Burn {
				Burn = true
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
