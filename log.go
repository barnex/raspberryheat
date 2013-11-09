package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logfile *os.File

const logfname = "temperature.log"

func doLog() {
	// make sure log file is open
	if logfile == nil {
		var err error
		logfile, err = os.OpenFile(logfname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("logging to", logfname)
		fmt.Fprintln(logfile)
	}

	// log
	fmt.Fprint(logfile, time.Now().Unix())
	for _, s := range sensor {
		fmt.Fprint(logfile, "\t", s.AvgTemp()) // also resets average temp
	}
	fmt.Fprintln(logfile)
}

func assureLogFile() {
}
