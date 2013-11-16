package main


type Room struct{
	sensor *Sensor
	Name string
}

func NewRoom(name, sensor string, led *GPIO)*Room{
	return &Room{Name: name, sensor: NewSensor(sensor, led)}
}

// GUILabel returns prefix + name, to be used as label for GUI element.
// E.g.: "somebutton_livingroom"
func (r *Room) GUILabel(prefix string) string {
	return prefix + "_" + r.Name
}

