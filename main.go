package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
)

var (
	Address = flag.String("address", "", "the address to host on")
	Port    = flag.Int("port", 8000, "the port to host on")
)

func getTemplate(name string) *template.Template {
	file := "views/" + name + ".html"
	return template.Must(template.New("").Funcs(nil).ParseFiles(file,
		"views/base.html"))
}

func main() {
	flag.Parse()
	endpoint := fmt.Sprintf("%v:%v", *Address, *Port)

	http.HandleFunc("/", index)
	http.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("assets"))))

	err := http.ListenAndServe(endpoint, nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func render(t *template.Template, w http.ResponseWriter,
	varmap map[string]interface{}) {
	err := t.ExecuteTemplate(w, "body", varmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
