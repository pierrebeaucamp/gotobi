package main

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	t := getTemplate("index")
	render(t, w, nil)
}
