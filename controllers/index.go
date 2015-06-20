package controllers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
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

func profile(w http.ResponseWriter, r *http.Request) {
	auth, _ := r.Cookie("GoTobiAuthToken")

	req, err := http.NewRequest("POST",
		"https://api.sandbox.paypal.com/v1/identity/openidconnect/userinfo?schema=openid",
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "Bearer "+auth.Value)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := getTemplate("profile")
	row := DB.QueryRow("SELECT * FROM projects WHERE id = $1",
		data["user_id"].(string))

	var name string
	var email string
	var bio string

	err = row.Scan(&name, &email, &bio)
	if err != nil {
		render(t, w, nil)
		return
	}

	varmap := map[string]interface{}{
		"name":  name,
		"email": email,
		"bio":   bio,
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
