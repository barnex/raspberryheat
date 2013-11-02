package raspberryheat

import (
	"fmt"
	"log"
	"os"
)

type GPIO struct {
	pin  int
	file *os.File
}

var (
	LED1   = &GPIO{pin: 17}
	LED2   = &GPIO{pin: 27}
	LED3   = &GPIO{pin: 22}
	LED4   = &GPIO{pin: 25}
	LED5   = &GPIO{pin: 24}
	RELAY1 = &GPIO{pin: 23}
	RELAY2 = &GPIO{pin: 18}
)

const GPIOPath = "/sys/class/gpio/"

func (p *GPIO) Export() {
	echo(GPIOPath+"export", p.pin)
}

func (p *GPIO) Unexport() {
	echo(GPIOPath+"unexport", p.pin)
}

func (p *GPIO) Direction(d string) {
	echo(fmt.Sprint(GPIOPath, "gpio", p.pin, "/direction"), d)
}

var (
	ON  = []byte("1")
	OFF = []byte("0")
)

func (p *GPIO) Set(value bool) {
	if p.file == nil {
		p.Export()
		p.Direction("out")
		fname := fmt.Sprint(GPIOPath, "gpio", p.pin, "/value")
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
