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
	Init()
	defer Cleanup()

}

func Init() {
	status.Export()
	status.Direction("out")
	status.Set(true)

	sensors := LsSensors()
	fmt.Println(sensors)
}

func Cleanup() {
	status.Set(false)
	status.Unexport()
}
