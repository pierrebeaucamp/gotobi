package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
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
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		//TODO
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		//TODO
		return
	}

	var info map[string]interface{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d, err := time.ParseDuration(info["expires_in"].(string) + "s")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "GoTobiAuthToken",
		Value:   info["access_token"].(string),
		Path:    "/",
		Expires: time.Now().Add(d),
	}
	http.SetCookie(w, &cookie)

	t := getTemplate("close")
	render(t, w, nil)
}
