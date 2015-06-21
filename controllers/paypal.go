package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ClientID  = "ASbG8fJ7mRm3KXfx0sWx1Pv2W0Nxy-L_aiM2x_Bqh7lrqME0KGGm94-BsaRIG-GvYsoEvZMteHuyU8_I"
	Secret    = "EJ0q4GOEJEpBNOCK2Bmord7mqHR-rs8Qr2Y-m2o-tCV00kszMdJiv3HUQDbh3CaJutsCc2IGjFB8Jf26"
	ReturnURI = "https://gotobi.herokuapp.com/login"
)

func getAuthToken() (string, error) {
	values := "grant_type=client_credentials"
	data := bytes.NewReader([]byte(values))

	req, err := http.NewRequest("POST",
		"https://api.paypal.com/v1/oauth2/token", data)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")
	req.SetBasicAuth(ClientID, Secret)

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

	var info map[string]interface{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return "", err
	}

	return info["access_token"].(string), nil

}

func PaypalLogin(w http.ResponseWriter, r *http.Request) {
	values := "client_id=" + ClientID + "&client_secret=" + Secret +
		"&grant_type=authorization_code&code=" + r.URL.Query().Get("code")
	data := bytes.NewReader([]byte(values))

	req, err := http.NewRequest("POST",
		"https://api.paypal.com/v1/identity/openidconnect/tokenservice",
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

func invoice(amount string, currency string, account string,
	email string) error {

	invoice := []byte(`{
		"merchant_info": {
			"email": "paypal_de@robocomobo.com"
		},
		"billing_info": [{
			"email": "` + email + `"
		}],
		"items": [{
			"name": "Donation",
			"quantity": 1,
			"unit_price": {
				"currency": "` + currency + `",
				"value": "` + amount + `"
			}
		}]
	}`)
	data := bytes.NewReader(invoice)

	req, err := http.NewRequest("POST",
		"https://api.paypal.com/v1/invoicing/invoices", data)
	if err != nil {
		return err
	}

	auth, err := getAuthToken()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	var info map[string]interface{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return err
	}

	if info["id"] == nil {
		fmt.Println(string(body))
		// graceful for now
		return nil
	}

	err = sendInvoice(info["id"].(string))
	if err != nil {
		return err
	}

	return nil
}

func sendInvoice(id string) error {
	req, err := http.NewRequest("POST",
		"https://api.paypal.com/v1/invoicing/invoices/"+
			id+"/send",
		nil)
	if err != nil {
		return err
	}

	fmt.Println("Invoice ID: " + id)

	auth, err := getAuthToken()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
