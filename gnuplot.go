package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func servePlot(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	todayLog := logFileName(now)
	yesterLog := logFileName(now.Add(-24 * time.Hour))

	cmd := `set xdata time; set xlabel "h"; set ylabel "deg C"; set timefmt "%s"; set format x "%H"; set term svg size 600,320 fsize 10; plot`
	for i, s := range sensor {
		if i != 0 {
			cmd += ","
		}
		using := i + 2
		cmd += fmt.Sprintf(`"<cat %v %v" u 1:%v w li title "%v"`, yesterLog, todayLog, using, s.Description())
	}
	cmd += `; set output;exit;`

	args := []string{"-e", cmd}
	//log.Println("gnuplot", args)
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
