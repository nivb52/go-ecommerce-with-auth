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

type TransactionData struct {
	FirstName       string
	LastName        string
	CardHolder      string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	Amount          int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

func (app *application) GetTransactionData(r *http.Request) (TransactionData, error) {
	var txData TransactionData
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return txData, err
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
	amount, detailedErr := convertAtoi(paymentAmount)
	if detailedErr != nil {
		app.errorLog.Println("failed to convert string to int:\n ", detailedErr)
		return txData, err
	}

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}
	pi, err := card.RetrivePaymentIntent(paymentIntent)
	if err != nil {
		app.infoLog.Println(" ::ERROR failed to retrive payment intent")
		app.errorLog.Println(err)
		return txData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.errorLog.Println("failed to get payment method, due to: \n", err)
		return txData, err
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear
	bankReturnCode := pi.Charges.Data[0].ID

	txData = TransactionData{
		FirstName:       customerFirstName,
		LastName:        customerLastName,
		CardHolder:      cardHolder,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		Amount:          amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     int(expiryMonth), // convert uint64 to int
		ExpiryYear:      int(expiryYear),  // convert uint64 to int
		BankReturnCode:  bankReturnCode,
	}

	return txData, nil
}

// PaymentSucceeded displays the receipt page
func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit PaymentSucceeded Handler")

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// get transaction data
	txData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	// read posted data
	widgetId, detailedErr := convertAtoi(r.Form.Get("product_id"))
	requestId := r.Form.Get("request_id")

	if detailedErr != nil {
		app.errorLog.Println("failed to convert string to int:\n ", detailedErr)
		return
	}

	// check if there is a re-submit of the form with the same requestId
	tables := models.NewModels(app.DB.DB)
	id, err := tables.DB.GetOrderByRequestId(requestId)
	if err != nil && id != 0 {
		app.infoLog.Println(" ::The form re-submited - abourting....")
		return
	}

	// create new customer
	customerID, customerSqlErr := app.SaveCustomer(txData.FirstName, txData.LastName, txData.Email)
	if err != nil {
		app.errorLog.Println("failed to save customer to DB due to:\n ", customerSqlErr)

	}
	app.infoLog.Println("customer has been created with ID: ", customerID)

	// create new transaction
	if customerSqlErr == nil {
		txn := models.Transaction{
			Amount:              txData.Amount,
			Currency:            txData.PaymentCurrency,
			LastFour:            txData.LastFour,
			ExpiryMonth:         txData.ExpiryMonth,
			ExpiryYear:          txData.ExpiryYear,
			BankReturnCode:      txData.BankReturnCode,
			PaymentIntent:       txData.PaymentIntentID,
			PaymentMethod:       txData.PaymentMethodID,
			TransactionStatusID: 2, //CLEARED,
		}

		txnID, err := app.SaveTxn(txn)
		if err != nil {
			app.errorLog.Println("failed to save transaction to DB due to:\n ", err)
			return
		}
		app.infoLog.Println("transaction has been created with ID: ", txnID)

		// create new order
		orderID, err := app.SaveOrder(
			widgetId,
			txnID,
			customerID,
			1, // CLEARED
			1, // WE allow only 1 quantity at this time
			txData.Amount,
			requestId,
		)
		if err != nil {
			app.errorLog.Println("failed to save order to DB due to:\n ", err)
		}
		app.infoLog.Println("order has been created with ID: ", orderID)
	}

	//data := make(map[string]interface{})
	//data["email"] = txData.Email
	//data["first_name"] = txData.FirstName
	//data["last_name"] = txData.LastName
	//
	//data["cardholder"] = txData.CardHolder
	//data["pi"] = txData.PaymentIntentID
	//data["pm"] = txData.PaymentMethodID
	//data["pa"] = txData.Amount
	//data["pc"] = txData.PaymentCurrency
	//
	//data["last_four"] = txData.LastFour
	//data["expiry_month"] = txData.ExpiryMonth
	//data["expiry_year"] = txData.ExpiryYear
	//data["bank_return_code"] = txData.BankReturnCode

	// write this data to session, and then redirect user to new page?
	app.Session.Put(r.Context(), "receipt", txData)
	http.Redirect(w, r, "/receipt", http.StatusOK)
	return
}

func (app *application) Receipt(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit Receipt Handler")

	txn := app.Session.Get(r.Context(), "receipt").(TransactionData)
	app.Session.Remove(r.Context(), "receipt")

	data := make(map[string]interface{})
	data["txn"] = txn

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
	err := app.renderTemplate(w, r, "single-widget", &templateData{
		Data: data,
	}, "stripe-js", "common-js")
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

func (app *application) SaveTxn(txn models.Transaction) (int, error) {

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
	requestId string,
) (int, error) {

	ordr := models.Order{
		WidgetID:      widgetId,
		TransactionID: transactionId,
		CustomerID:    customerId,
		StatusID:      statusId,
		Quantity:      quantity,
		Amount:        amount,
		RequestID:     requestId,
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
