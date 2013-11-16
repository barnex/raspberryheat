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
		doc.SetValue("burn", Burn)
		for _, r := range rooms {
			doc.SetValue(r.GUILabel("readout"), fmt.Sprintf("%.1f", r.sensor.Temp()))
			doc.SetValue(r.GUILabel("error"), r.sensor.Error())
			doc.SetValue(r.GUILabel("settemp"), r.SetTemp)
			doc.SetValue(r.GUILabel("burn"), r.Burn)
		}
	})

for _, r := range rooms {
	r := r
	label := r.GUILabel("settemp")
	doc.OnEvent(label,func(){ r.SetTemp = doc.Value(label).(float64)})
}

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
	brander: {{$.Span "burn"}} 

	<hr/>
	
	{{ range $.Data }}

		<h2>{{.Name}}</h2>
		<span style="font-size:2em; font-weight:bold">
			{{.GUILabel "readout" | $.Span }} <sup>o</sup>C 
		</span><br/>
		<span style="font-weight:bold; color:red"> {{.GUILabel "error" | $.Span}} </span> <br/>

		set: {{.GUILabel "settemp" | $.NumBox}} <sup>o</sup>C <br/>
		brander: {{.GUILabel "burn" | $.Span}} <br/>

	<hr/>

	{{ end }}

	<img src="plot">

	<hr/>

	<p> {{.ErrorBox}} </p>

</body>
</html>
`
