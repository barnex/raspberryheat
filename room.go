package main

const (
	MINTEMP        = 6   // no room should be colder than this (°C), ever.
	DEFAULT_WINDOW = 0.5 // Schmidt trigger temperature window (°C)
)

type Room struct {
	sensor  *Sensor
	Name    string
	Burn    bool
	SetTemp float64
	Schmidt float64
}

func NewRoom(name, sensor string, led *GPIO) *Room {
	return &Room{Name: name, sensor: NewSensor(sensor, led),
		SetTemp: MINTEMP, Schmidt: DEFAULT_WINDOW}
}

func(r*Room)Update(){
			r.sensor.Update()
			if r.sensor.Temp() > r.SetTemp+r.Schmidt/2 {
				r.Burn = false
			}
			if r.sensor.Temp() < r.SetTemp-r.Schmidt/2 {
				r.Burn = true
			}
}

// GUILabel returns prefix + name, to be used as label for GUI element.
// E.g.: "somebutton_livingroom"
func (r *Room) GUILabel(prefix string) string {
	return prefix + "_" + r.Name
}
