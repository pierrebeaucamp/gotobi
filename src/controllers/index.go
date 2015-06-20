package controllers

import (
	"html/template"
	"net/http"
)

func getTemplate(name string) *template.Template {
	file := "src/views/" + name + ".html"
	return template.Must(template.New("").Funcs(nil).ParseFiles(file,
		"src/views/base.html"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := getTemplate("index")
	render(t, w, nil)
}

func render(t *template.Template, w http.ResponseWriter,
	varmap map[string]interface{}) {
	err := t.ExecuteTemplate(w, "body", varmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
