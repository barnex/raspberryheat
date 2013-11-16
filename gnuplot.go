package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var dst = 1 * time.Hour

func servePlot(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	todayLog := logFileName(now)
	yesterLog := logFileName(now.Add(-24 * time.Hour))

	cmd := `set grid; set xdata time; set xlabel "h"; set ylabel "deg C"; set timefmt "%s"; set format x "%H"; set term svg size 600,320 fsize 10; plot`
	for i, r := range rooms {
		if i != 0 {
			cmd += ","
		}
		usingX := fmt.Sprint("($1+", dst.Seconds(), ")")
		usingY := i + 2
		cmd += fmt.Sprintf(`"<cat %v %v" u %v:%v w li title "%v"`, yesterLog, todayLog, usingX, usingY, r.Name)
	}
	cmd += `; set output;exit;`

	args := []string{"-e", cmd}
	log.Println("gnuplot", args)
	out, err := exec.Command("gnuplot", args...).Output()

	if err != nil {
		log.Println(err)
		log.Println(string(out))
		msg := string(out) + "\n" + err.Error()
		http.Error(w, msg, 400)
		return
	} else {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-control", "No-Cache")
		w.Write(out)
	}
}
