package repository

import (
	"database/sql"
	"fmt"
)

const (
	userTable         = "users"
	itemsTable        = "items"
	salesTable        = "sales"
	transactionsTable = "transactions"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return db, nil
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil
}
