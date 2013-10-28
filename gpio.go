package raspberryheat

import (
	"fmt"
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

func echo(dest string, msg interface{}) {
	f, err := os.OpenFile(dest, os.O_WRONLY, 0666)
	check(err)
	defer f.Close()
	checkIO(fmt.Fprint(f, msg))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkIO(n int, err error) {
	check(err)
}
