package raspberryheat

import (
	"fmt"
	"log"
	"os"
)

type Gpio struct {
	pin  int
	file *os.File
}

func GPIO(pin int) *Gpio {
	return &Gpio{pin: pin}
}

const GpioPath = "/sys/class/gpio/"


func (p *Gpio) Export() {
	echo(GpioPath+"export", p.pin)
}

func (p *Gpio) Unexport() {
	echo(GpioPath+"unexport", p.pin)
}

func (p *Gpio) Direction(d string) {
	echo(fmt.Sprint(GpioPath, "gpio", p.pin, "/direction"), d)
}

var (
	ON  = []byte("1")
	OFF = []byte("0")
)

func (p *Gpio) Set(value bool) {
	if p.file == nil {
		fname := fmt.Sprint(GpioPath, "gpio", p.pin, "/value")
		f, err := os.OpenFile(fname, os.O_WRONLY, 0666)
		if err != nil {
			Log(err)
			return
		} else {
			p.file = f
		}
	}
	if value {
		checkIO(p.file.Write(ON))
	} else {
		checkIO(p.file.Write(OFF))
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
