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
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return db, nil
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil
}
