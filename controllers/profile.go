package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

func getUserID(auth string) (string, error) {
	req, err := http.NewRequest("POST",
		"https://api.sandbox.paypal.com/v1/identity/openidconnect/userinfo?schema=openid",
		nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+auth)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data["user_id"].(string), nil
}

func profile(w http.ResponseWriter, r *http.Request) {
	auth, _ := r.Cookie("GoTobiAuthToken")
	user_id, err := getUserID(auth.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := getTemplate("profile")
	row := DB.QueryRow(`SELECT name, email, bio FROM projects WHERE id = $1`,
		user_id)

	var name string
	var email string
	var bio string

	err = row.Scan(&name, &email, &bio)
	if err != nil {
		fmt.Println(err.Error())
		render(t, w, nil)
		return
	}

	varmap := map[string]interface{}{
		"name":   name,
		"email":  email,
		"bio":    bio,
		"filled": true,
	}

	render(t, w, varmap)
}

func Submit(w http.ResponseWriter, r *http.Request) {
	auth, err := r.Cookie("GoTobiAuthToken")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}

	user_id, err := getUserID(auth.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("filled") != "true" {
		_, err = DB.Exec(`INSERT INTO projects (id, name, bio, email) 
			VALUES ($1, $2, $3, $4)`, user_id, r.FormValue("title"),
			r.FormValue("bio"), r.FormValue("email"))
	} else {
		_, err = DB.Exec(`UPDATE projects SET (name, bio, email) = 
			($2, $3, $4) WHERE id = $1`, user_id, r.FormValue("title"),
			r.FormValue("bio"), r.FormValue("email"))
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
