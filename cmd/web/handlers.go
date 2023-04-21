package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) Liveness(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal("{liveness: true, code: 200}")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
	w.Write(jsonBytes)
}

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit VirtualTerminal Handler")
	// stringMap := make(map[string]string)
	// stringMap["publishable_key"] = app.config.stripe.key

	if err := app.renderTemplate(w, r, "terminal", &templateData{
		// StringMap: stringMap,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// PaymentSucceeded displays the receipt page
func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit PaymentSucceeded Handler")

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// read posted data
	cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency

	// should write this data to session, and then redirect user to new page?

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}
