package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ClientID  = "ASTQe8TCWn2Ulfw4nEmOwCWBOJ5efoRdr6ajtVZYwJGdGcDpfOEG509CwKkRlAV5R-MxtpWhcV2X5GpZ"
	Secret    = "EAS3qVHJG9MBg0Lme_QeTX_-ex7e2CpqE2pKLhGISHkid1eAXjgY1RdBL2Y_0WY-GQOV3Dj-LKRNnG1p"
	ReturnURI = "http://localhost:8000/login/"
)

func PaypalLogin(w http.ResponseWriter, r *http.Request) {
	values := "client_id=" + ClientID + "&client_secret=" + Secret +
		"&grant_type=authorization_code&code=" + r.URL.Query().Get("code")
	data := bytes.NewReader([]byte(values))

	req, err := http.NewRequest("POST",
		"https://api.sandbox.paypal.com/v1/identity/openidconnect/tokenservice",
		data)
	if err != nil {
		//TODO
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		//TODO
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		//TODO
	}

	fmt.Fprint(w, string(body))
}
