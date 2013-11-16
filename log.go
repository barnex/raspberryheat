package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	logFile    *os.File
	logFileDay int
)

func doLog() {
	// make sure log file is open
	if logFile == nil {
		now := time.Now()
		fname := logFileName(now)
		logFileDay = now.Day()
		_ = os.Mkdir(fmt.Sprint(now.Year()), 0777)
		var err error
		logFile, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("logging to", fname)
		fmt.Fprintln(logFile)
	}

	now := time.Now()

	if now.Day() != logFileDay {
		logFile.Close()
		logFile = nil
		return
	}

	// log
	fmt.Fprint(logFile, now.Unix())
	for _, r := range rooms {
		fmt.Fprintf(logFile, "\t%.3f", r.sensor.AvgTemp()) // also resets average temp
	}
	fmt.Fprintln(logFile)
}

func logFileName(now time.Time) string {
	return fmt.Sprint(now.Year(), "/", now.Month(), "-", now.Day(), ".log")
}
