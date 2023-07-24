package main

import (
	"encoding/json"
	"go-ecommerce-with-auth/internal/cards"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type stripPayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) Liveness(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal("{liveness: true, code: 200}")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	} else if amount <= 0 {
		app.infoLog.Println("payment is negative or zero")
		w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		j := jsonResponse{
			OK:      true,
			Message: "Payment is too low",
			Content: "Payment of " + payload.Amount + " is too low.",
		}

		out, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			app.errorLog.Println("parse json failed: ", err)
		}
		w.Write(out)
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true
	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		app.infoLog.Println("Card Charged")
		type CheckoutData struct {
			ClientSecret string `json:"client_secret"`
		}
		data := CheckoutData{
			ClientSecret: pi.ClientSecret,
		}

		w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)

	} else {
		app.infoLog.Println("Payment failed")
		j := jsonResponse{
			OK:      true,
			Message: msg,
			Content: "",
		}

		out, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			app.errorLog.Println("parse json failed: ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		// TODO: make globbal not found fnuc
		app.errorLog.Println(err)
		jsonBytes, _ := json.Marshal("{message: 'Widget Not Found', code: 404}")
		w.Header().Set("Content-type", "application/json")
		w.Write(jsonBytes)
	}

	out, jsonErr := json.Marshal(widget)
	if jsonErr != nil {
		app.errorLog.Println(jsonErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
