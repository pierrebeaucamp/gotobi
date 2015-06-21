package controllers

import (
	"fmt"
	"net/http"
)

func InboundMail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("from") == "" && r.FormValue("to") == "" {
		fmt.Fprint(w, "This is not an Email")
	}

	fmt.Println(r.FormValue("from"))
	fmt.Println(r.FormValue("subject"))
	fmt.Println(r.FormValue("to"))
}
