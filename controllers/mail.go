package controllers

import (
	"fmt"
	"net/http"
)

func InboundMail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("from") == "" && r.FormValue("subject") == "" {
		fmt.Fprint(w, "This is not an Email")
	}

	fmt.Fprint(w, r.FormValue("from")+"\n\n"+r.FormValue("subject"))
}
