package raspberryheat

import (
	"fmt"
	"log"
	"os"
)

const GpioPath = "/sys/class/gpio/"

type GPIO int

func (pin GPIO) String() string {
	return fmt.Sprint("gpio", int(pin))
}

func (pin GPIO) Export() {
	echo(GpioPath+"export", int(pin))
}

func (pin GPIO) Unexport() {
	echo(GpioPath+"unexport", int(pin))
}

func (pin GPIO) Direction(d string) {
	echo(GpioPath+pin.String()+"/direction", d)
}

func (pin GPIO) Set(value bool) {
	ctl := GpioPath + pin.String() + "/value"
	if value {
		echo(ctl, 1)
	} else {
		echo(ctl, 0)
	}
}

func echo(dest string, msg interface{}) {
	f, err := os.OpenFile(dest, os.O_WRONLY, 0666)
	if err != nil {
		Log(err)
		return
	}
	defer f.Close()
	checkIO(fmt.Fprint(f, msg))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Log(err error) {
	if err != nil {
		log.Println(err)
	}
}

func checkIO(n int, err error) {
	Log(err)
}
