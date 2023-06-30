package models

import (
	"context"
	"database/sql"
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
	CreatedAt      time.Time `json:"-"`
	UpdateAt       time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
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

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget
	query := "SELECT id, name, price, description, inventory_level FROM widgets WHERE id = ?"
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&widget.ID, &widget.Name, &widget.Price, &widget.Description, &widget.InventoryLevel)
	if err != nil {
		return widget, err
	}

	return widget, nil
}
