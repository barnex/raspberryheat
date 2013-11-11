package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func servePlot(w http.ResponseWriter, r *http.Request) {
	//url := r.URL.Path[len("/plot/"):]

	now := time.Now()
	todayLog := logFileName(now)
	yesterLog := logFileName(now.Add(-24 * time.Hour))

	cmd := fmt.Sprintf(`set key off; set xdata time; set timefmt "%%s"; set format x "%%H"; set term svg size 600,320 fsize 10; plot "%v" u 1:2 w li, "%v" u 1:2 w li; set output;exit;`, yesterLog, todayLog)
	args := []string{"-e", cmd}
	log.Println("gnuplot", args)
	out, err := exec.Command("gnuplot", args...).CombinedOutput()
	if err != nil {
		msg := string(out) + "\n" + err.Error()
		http.Error(w, msg, 400)
		return
	} else {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(out)
	}
}
