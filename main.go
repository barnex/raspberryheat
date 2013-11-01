//+build ignore

package main

import (
	. "."
	"fmt"
	"log"
	"time"
)

var(
	status = LED1
)

func main() {
	defer Cleanup()

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
		time.Sleep(10*time.Millisecond)
	}

}

func Cleanup() {
	status.Set(false)
	status.Unexport()
}
