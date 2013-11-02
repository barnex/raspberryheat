package main

import (
	"fmt"
	"net/http"
)

func StartHTTP() {
	http.HandleFunc("/", rootHandler)
	check(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, getTemp())
	blink(httpLED)
}
