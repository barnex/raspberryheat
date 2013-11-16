package main

import (
	"fmt"
	"github.com/barnex/gui"
	"net/http"
	"time"
)

func StartHTTP() {
	doc := gui.NewDoc(templ, rooms)

	doc.OnRefresh(func() {
		doc.SetValue("time", time.Now().Format(time.ANSIC))
		for _, r := range rooms {
			doc.SetValue(r.GUILabel("readout"), fmt.Sprintf("%.1f", r.sensor.Temp()))
			doc.SetValue(r.GUILabel("error"), r.sensor.Error())
		}
	})

	http.Handle("/", doc)
	http.HandleFunc("/plot/", servePlot)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

const templ = `

<html>

<head>
	<meta http-equiv="refresh" content="120">
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

		{{.Name}}
		{{.GUILabel "temp" | $.TextBox}} <sup>o</sup>C
		van {{.GUILabel "start" | $.TextBox}}
		tot {{.GUILabel "stop" | $.TextBox}} <br/>

		<span style="font-size:2em; font-weight:bold">
			{{.GUILabel "readout" | $.Span }} <sup>o</sup>C 
		</span><br/>
		<span style="font-weight:bold; color:red"> {{.GUILabel "error" | $.Span}} </span> <br/>

	<hr/>

	{{ end }}

	<img src="plot">

	<hr/>

	<p> {{.ErrorBox}} </p>

</body>
</html>
`
