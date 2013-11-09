package main

import (
	"net/http"
	"os/exec"
)

func servePlot(w http.ResponseWriter, r *http.Request) {
	//url := r.URL.Path[len("/plot/"):]

	cmd := "gnuplot"
	args := []string{"-e", `set xdata time; set timefmt "%s"; set format x "%H"; set term svg size 480,320 fsize 10; plot "temperature.log" u 1:2 w li; set output;exit;`}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	} else {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(out)
	}
}
