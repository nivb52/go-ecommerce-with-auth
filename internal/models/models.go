package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
)

// DBModel is type for db connection values
type DBModel struct {
	DB *sql.DB
}

// Wrapper to all models
type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"-"`
	UpdateAt       time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdateAt      time.Time `json:"-"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

type TransactionStatuses struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

type Transaction struct {
	ID                  int    `json:"id"`
	Amount              int    `json:"amount"`
	Currency            string `json:"currency"`
	LastFour            string `json:"last_four"`
	ExpiryMonth         int    `json:"expiry_month"`
	ExpiryYear          int    `json:"expiry_year"`
	BankReturnCode      string `json:"bank_return_code"`
	TransactionStatusID int    `json:"transaction_statuses_id"`

	CreatedAt time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	CreatedAt time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

var errorLog = log.New(os.Stdout, ":: ERROR SQL:\t", log.Ltime)
var sillyLog = log.New(os.Stdout, ":: Silly:\t", log.Ltime)

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget
	query := `SELECT 
		id, name, price, description, inventory_level, coalesce(image, '') 
		FROM widgets 
		WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Price,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Image)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

// InsertTransaction insert new txn, and return its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	quary := `
		INSERT INTO transactions 
		(amount, currency, last_four, bank_return_code, transaction_status_id)
	VALUES (?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, quary,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
	)
	if err != nil {
		errorLog.Println("Transaction execution failed due:\n", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		errorLog.Println("Transaction executed but failed to retrive last inserted id, err: \n", err)
		return 0, err
	}

	return int(id), nil
}

// InsertOrder insert new order, and return its id
func (m *DBModel) InsertOrder(ordr Order) (int, error) {
	quary := `
		INSERT INTO orders 
		(widget_id, transaction_id, customer_id, 
			status_id, quantity, amount
		)
		VALUES (?, ?, ?, 
				?, ?, ?
			)
	`
	result, err := m.insertQuery(quary, "order",
		ordr.WidgetID,
		ordr.TransactionID,
		ordr.CustomerID,
		ordr.StatusID,
		ordr.Quantity,
		ordr.Amount,
	)

	if err != nil {
		sillyLog.Println("Order execution failed due")
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		errorLog.Println("Order executed but failed to retrive last inserted id, err: \n", err)
		return 0, err
	}

	return int(id), nil
}

// InsertCustomer insert new customer, and return its id
func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	quary := `
		INSERT INTO customers 
		(first_name, last_name, email)
	VALUES (?, ?, ?)
	`
	result, err := m.insertQuery(quary, "customer",
		c.FirstName,
		c.LastName,
		c.Email,
	)

	if err != nil {
		sillyLog.Println("Customer execution failed")
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		errorLog.Println("Customer executed but failed to retrive last inserted id, err: \n", err)
		return 0, err
	}

	return int(id), nil
}

// insertQuary
func (m *DBModel) insertQuery(query, logAnnouncment string, data ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, data...)
	if err != nil {
		errorLog.Println(logAnnouncment, " execution failed due:\n", err)
		return result, err
	}

	return result, nil
}
