package controllers

import (
	"database/sql"
	"html/template"
	"net/http"
)

var (
	DB *sql.DB
)

func getTemplate(name string) *template.Template {
	file := "src/html/" + name + ".html"
	return template.Must(template.New("").Funcs(nil).ParseFiles(file,
		"src/html/base.html"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("GoTobiAuthToken")
	if err == nil {
		profile(w, r)
		return
	}

	t := getTemplate("index")

	varmap := map[string]interface{}{
		"clientID":  ClientID,
		"returnURI": ReturnURI,
	}

	render(t, w, varmap)
}

func render(t *template.Template, w http.ResponseWriter,
	varmap map[string]interface{}) {
	err := t.ExecuteTemplate(w, "body", varmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
