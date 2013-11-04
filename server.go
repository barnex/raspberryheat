package main

import (
	"github.com/barnex/gui"
	"net/http"
	"time"
)

func StartHTTP() {
	data := []int{0, 1} // room numbers
	doc := gui.NewDoc(templ, data)

	doc.OnRefresh(func() {
		doc.SetValue("time", time.Now().Format(time.ANSIC))
		doc.SetValue("temp0", sensor[0].Temp())
		doc.SetValue("temp1", sensor[1].Temp())
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

	<p> {{.ErrorBox}} </p>
	{{.Span "time" "--"}}  <br/>

	<hr/>

	living <p style="font-size:2em; font-weight:bold">{{.Span "temp0" "--"}} <sup>o</sup>C <p> <br/>
	kindjes <p style="font-size:2em; font-weight:bold">{{.Span "temp1" "--"}} <sup>o</sup>C <p> <br/>

	<hr/>

</body>
</html>
`
