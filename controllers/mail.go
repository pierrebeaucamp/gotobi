package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func getAccount(in string) string {
	return strings.Split(in, "@")[0]
}

func getAmount(in string) (string, string) {
	regCur := regexp.MustCompile(`\p{Sc}|AUD|BRL|CAD|CZK|DKK|EUR|HKD|HUF|ILS|JPY|MYR|MXN|TWD|NZD|NOK|PHP|PLN|GBP|SGD|SEK|CHF|THB|TRY|USD`)
	regAm := regexp.MustCompile(`[0-9]*$`)

	currency := regCur.FindStringSubmatch(in)[0]
	if currency == "" {
		// Default
		currency = "EUR"
	}

	switch currency {
	case "$":
		currency = "USD"
	case "\u20AC":
		currency = "EUR"
	}

	// 12.34 EUR
	amount := regCur.Split(in, -1)[0]
	if amount == "" {
		// EUR 12.34
		amount = regCur.Split(in, -1)[1]
	}

	amount = regAm.FindStringSubmatch(amount)[0]
	if amount == "" {
		return "0", currency
	}

	return amount, currency
}

func InboundMail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("from") == "" && r.FormValue("to") == "" {
		fmt.Fprint(w, "This is not an Email")
	}

	amount, currency := getAmount(r.FormValue("subject"))
	fmt.Println("Got amount: " + amount + " and currency: " + currency)

	if amount == "0" {
		// Respond 200 so sendgrid doesn't try over and over again
		fmt.Fprintf(w, "200")
		return
	}

	fmt.Println(r.PostFormValue("to"))
	account := getAccount(r.PostFormValue("to"))
	fmt.Println(account)
	if account == "" {
		fmt.Fprint(w, "200")
		return
	}

	err := invoice(amount, currency, account, r.FormValue("From"))
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
}
