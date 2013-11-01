package raspberryheat

import (
	"fmt"
	"log"
	"os"
)

type gpio struct {
	pin  int
	file *os.File
}

var (
	LED1   = &gpio{pin: 17}
	LED2   = &gpio{pin: 27}
	LED3   = &gpio{pin: 22}
	LED4   = &gpio{pin: 25}
	LED5   = &gpio{pin: 24}
	RELAY1 = &gpio{pin: 23}
	RELAY2 = &gpio{pin: 18}
)

const GpioPath = "/sys/class/gpio/"

func (p *gpio) Export() {
	echo(GpioPath+"export", p.pin)
}

func (p *gpio) Unexport() {
	echo(GpioPath+"unexport", p.pin)
}

func (p *gpio) Direction(d string) {
	echo(fmt.Sprint(GpioPath, "gpio", p.pin, "/direction"), d)
}

var (
	ON  = []byte("1")
	OFF = []byte("0")
)

func (p *gpio) Set(value bool) {
	if p.file == nil {
		p.Export()
		p.Direction("out")
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
