//+build ignore

package main

import (
	. "."
	"fmt"
)

const (
	status = GPIO(24)
)

func main() {
	defer Cleanup()

	status.Export()
	status.Direction("out")
	status.Set(true)


	sensors := LsSensors()
	fmt.Println(sensors)

	for{
		status.Set(false)
		for _,s := range sensors{
			fmt.Println(s, s.Read())
		}
		status.Set(true)
		fmt.Println()

	}

}

func Cleanup() {
	status.Set(false)
	status.Unexport()
}
