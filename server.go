package main

import (
	"net/http"
	"text/template"
)

func StartHTTP() {
	var data Dummy
	templ := template.Must(template.New("root").Parse(templText))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, data)
		blink(httpLED)
	})
	check(http.ListenAndServe(":8080", nil))
}

type Dummy int

const templText = `
	
`
