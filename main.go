//+build ignore

package main

import (
	. "."
	"fmt"
	"log"
)

const (
	LED1 = GPIO(17)
	LED2 = GPIO(27)
	LED3 = GPIO(22)
	LED4 = GPIO(25)
	LED5 = GPIO(24)
	RELAY1 = GPIO(23)
	RELAY2 = GPIO(18)
	status = RELAY2
)

func main() {
	defer Cleanup()

	status.Export()
	status.Direction("out")
	status.Set(true)


	sensors := LsSensors()
	log.Println(sensors)

	for{
		status.Set(false)
		for _,s := range sensors{
			t, err := s.Read()
			if err != nil{
				log.Println(err)
				continue
			}else{
				fmt.Println(t)
			}
		}
		status.Set(true)
	}

}

func Cleanup() {
	status.Set(false)
	status.Unexport()
}
