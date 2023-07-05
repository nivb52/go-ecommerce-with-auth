package main

import (
	"encoding/json"
	"go-ecommerce-with-auth/internal/cards"
	"go-ecommerce-with-auth/internal/models"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) Liveness(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal("{liveness: true, code: 200}")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+os.Getenv("GOSTRIPE_PORT"))
	w.Write(jsonBytes)
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit VirtualTerminal Handler")
	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
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
	customerFirstName := r.Form.Get("first_name")
	customerLastName := r.Form.Get("last_name")
	cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	_, detailedErr := convertAtoi(r.Form.Get("product_id"))
	if detailedErr != nil {
		app.errorLog.Println("::ERROR : failed to convert string to int:\n ", detailedErr)
		return
	}

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}
	pi, err := card.RetrivePaymentIntent(paymentIntent)
	if err != nil {
		app.infoLog.Println(" ::ERROR failed to retrive payment intent")
		app.errorLog.Println(err)
		return
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.errorLog.Println("failed to get payment method, due to: \n", err)
		return
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear
	bankReturnCode := pi.Charges.Data[0].ID

	// create new customer
	customerID, customerSqlErr := app.SaveCustomer(customerFirstName, customerLastName, email)
	if err != nil {
		app.errorLog.Println("failed to save customer to DB due to:\n ", customerSqlErr)

	}
	app.infoLog.Println("customer has been created with ID: ", customerID)

	amount, detailedErr := convertAtoi(paymentAmount)
	if detailedErr != nil {
		app.errorLog.Println("failed to convert string to int:\n ", detailedErr)
		return
	}

	data := make(map[string]interface{})
	data["email"] = email
	data["first_name"] = customerFirstName
	data["last_name"] = customerLastName

	data["cardholder"] = cardHolder
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = amount
	data["pc"] = paymentCurrency

	data["last_four"] = lastFour
	data["expiry_month"] = expiryMonth
	data["expiry_year"] = expiryYear
	data["bank_return_code"] = bankReturnCode

	// should write this data to session, and then redirect user to new page?

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) WidgetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)
	db := app.DB
	widget, dbErr := db.GetWidget(widgetID)
	if dbErr != nil {
		app.errorLog.Println(dbErr)
		jsonBytes, _ := json.Marshal("{message: 'Widget Not Found', code: 404}")
		w.Header().Set("Content-type", "application/json")
		w.Write(jsonBytes)
		return
	}

	data := make(map[string]interface{})
	data["widget"] = widget
	app.debugLog.Println("widget.Price ", widget.Price)
	err := app.renderTemplate(w, r, "buy-once", &templateData{
		Data: data,
	}, "stripe-js")
	if err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) SaveCustomer(firstName string, lastName string, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	tables := models.NewModels(app.DB.DB)
	id, err := tables.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *application) SaveTxn(amount int,
	currency string,
	lastFour string,
	expiryMonth int,
	expiryYear int,
	bankReturnCode string,
	transactionStatusId int,
) (int, error) {

	txn := models.Transaction{
		Amount:              amount,
		Currency:            currency,
		LastFour:            lastFour,
		ExpiryMonth:         expiryMonth,
		ExpiryYear:          expiryYear,
		BankReturnCode:      bankReturnCode,
		TransactionStatusID: transactionStatusId,
	}

	tables := models.NewModels(app.DB.DB)
	id, err := tables.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *application) SaveOrder(widgetId int,
	transactionId int,
	customerId int,
	statusId int,
	quantity int,
	amount int,
) (int, error) {

	ordr := models.Order{
		WidgetID:      widgetId,
		TransactionID: transactionId,
		CustomerID:    customerId,
		StatusID:      statusId,
		Quantity:      quantity,
		Amount:        amount,
	}

	tables := models.NewModels(app.DB.DB)
	id, err := tables.DB.InsertOrder(ordr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func convertAtoi(n string) (int, error) {
	res, err := strconv.Atoi(n)
	if err != nil {
		detailedErr := &strconv.NumError{
			Num:  n,
			Err:  err,
			Func: "Atoi",
		}
		return 0, detailedErr
	}
	return res, nil
}
