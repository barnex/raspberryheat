package main

import (
	"math"
)

const (
	MINTEMP        = 6   // no room should be colder than this (°C), ever.
	MAXTEMP        = 25   // no room should be hotter than this (°C), ever.
	DEFAULT_WINDOW = 0.5 // Schmidt trigger temperature window (°C)
)

type Room struct {
	sensor  *Sensor
	Name    string  // E.g. livingroom
	Burn    bool    // This room currently asks heat?
	SetTemp float64 // Current desired temperture
	Schmidt float64 // Schmidt-trigger window on SetTemp
}

func NewRoom(name, sensor string, led *GPIO) *Room {
	return &Room{Name: name, sensor: NewSensor(sensor, led),
		SetTemp: MINTEMP, Schmidt: DEFAULT_WINDOW}
}

func (r*Room)SetSchmidt(dT float64){
	if dT < 0.1{
		dT = 0.1
	}
	r.Schmidt = dT
}


func (r*Room)SetSetTemp(T float64){
	if T < MINTEMP{
	T = MINTEMP
	}
	if T > MAXTEMP{
		T= MAXTEMP
	}
	r.SetTemp = T
}

// Update burn state etc.
func (r *Room) UpdateBurn() {
	temp := r.sensor.Temp()
	switch {
	case temp > r.SetTemp+r.Schmidt:
		r.Burn = false
	case temp < r.SetTemp-r.Schmidt:
		r.Burn = true
	case math.IsNaN(temp): // sensor disconnected
		r.Burn = false
	}
}

func (r*Room)Overheat()bool{
	return r.sensor.temp > r.SetTemp+r.Schmidt
}

// GUILabel returns prefix + name, to be used as label for GUI element.
// E.g.: "somebutton_livingroom"
func (r *Room) GUILabel(prefix string) string {
	return prefix + "_" + r.Name
}
