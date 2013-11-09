package main

import (
	"fmt"
	"github.com/barnex/gui"
	"net/http"
	"time"
)

func StartHTTP() {
	doc := gui.NewDoc(templ, sensor)

	doc.OnRefresh(func() {
		doc.SetValue("time", time.Now().Format(time.ANSIC))
		for _, s := range sensor {
			doc.SetValue(s.Label("readout"), fmt.Sprintf("%.1f", s.Temp()))
			doc.SetValue(s.Label("error"), s.Error())
		}
	})

	http.Handle("/", doc)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

const templ = `

<html>

<head>
	<style type="text/css">
		body      { margin: 20px; font-family: Ubuntu, Arial, sans-serif; }
		hr        { border-style: none; border-top: 1px solid #CCCCCC; }
		.ErrorBox { color: red; font-weight: bold; } 
		.TextBox  { border:solid; border-color:#BBBBBB; border-width:1px; padding-left:4px;}
	</style>
	{{.JS}}
</head>

<body>

	{{.Span "time"}}  <br/>

	<hr/>
		master:
		{{$.Button "masterOn"   "<b>ON </b>"}} 
		{{$.Button "masterOff"  "<b>OFF</b>"}} 
		{{$.Button "masterAuto" "<b>Auto</b>"}} 

	<hr/>
	
	{{ range $.Data }}

		{{.Description}}
		{{.Label "temp" | $.TextBox}} <sup>o</sup>C
		van {{.Label "start" | $.TextBox}}
		tot {{.Label "stop" | $.TextBox}} <br/>

		<span style="font-size:2em; font-weight:bold">
			{{.Label "readout" | $.Span }} <sup>o</sup>C 
		</span><br/>
		<span style="font-weight:bold; color:red"> {{.Label "error" | $.Span}} </span> <br/>

	{{ end }}

	<hr/>

	<p> {{.ErrorBox}} </p>

</body>
</html>
`
