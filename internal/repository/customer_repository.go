package repository

import (
	"database/sql"
	"encoding/json"
	"frappuchino/internal/models"
	"log/slog"

	_ "github.com/lib/pq"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) Close() error {
	return r.db.Close()
}

func (r *CustomerRepository) IndentCustomerID(customerName string, instructions json.RawMessage) (int, error) {
	var customerID int
	query := `
		SELECT id
		FROM customers
		WHERE name = $1
	`
	err := r.db.QueryRow(query, customerName).Scan(&customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			insertQuery := `
				INSERT INTO customers (name, email, preferences)
				VALUES ($1, $2, $3)
				RETURNING id
			`
			customer, err := models.NewCustomer(customerName, instructions)
			if err != nil {
				slog.Error("Repository error from Ident Customer ID: invalid input data", "customer name", customerName, "error", err)
				return 0, err
			}

			if err := r.db.QueryRow(insertQuery, customer.Name, customer.Email, customer.Preferences).Scan(&customerID); err != nil {
				slog.Error("Repository error from Ident Customer ID: failed to insert into table", "customer", customer, "error", err)
				return 0, err
			}
		} else {
			slog.Error("Repository error from Ident Customer ID: failed to select from table", "customer name", customerName, "error", err)
			return 0, err
		}
	}

	slog.Info("Repository info: ident customer ID successfully", "customer ID", customerID)
	return customerID, nil
}
