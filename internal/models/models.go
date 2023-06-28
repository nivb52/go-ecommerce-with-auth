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

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget
	query := "SELECT id, name, price FROM widgets WHERE id = ?"
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&widget.ID, &widget.Name, &widget.Price)
	if err != nil {
		return widget, err
	}

	return widget, nil
}
