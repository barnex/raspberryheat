package main

import (
	"time"
)

const logPeriod = 2 * time.Minute
var HEATER = RELAY2

var rooms = []*Room{
	NewRoom("living", "28-0000050ad012", LED1),
	NewRoom("kindjes", "28-0000050b07f7", LED2)}


func main() {
	go StartHTTP()

	lastLog := time.Now()
	for {
		for _, r := range rooms {
			r.sensor.Update()
		}
		if time.Since(lastLog) > logPeriod {
			doLog()
			lastLog = time.Now()
		}

		if rooms[0].sensor.Temp() > 18.2{
			HEATER.Set(false)
		}
		if rooms[0].sensor.Temp() < 17.8{
			HEATER.Set(true)
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
