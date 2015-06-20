package gotobi

import (
	"fmt"
	"net/http"
)

func init() {
	r := http.NewServeMux()
	r.HandleFunc("/", index)
	http.Handle("/", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
