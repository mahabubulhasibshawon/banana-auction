package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"banana-auction/config"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func InitDB(cfg *config.Config) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL CHECK (role IN ('seller', 'buyer'))
		);
		CREATE TABLE IF NOT EXISTS lots (
			id SERIAL PRIMARY KEY,
			seller_id INTEGER REFERENCES users(id),
			cultivar TEXT NOT NULL,
			planted_country TEXT NOT NULL,
			harvest_date TEXT NOT NULL,
			total_weight_kg INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS auctions (
			id SERIAL PRIMARY KEY,
			lot_id INTEGER REFERENCES lots(id),
			start_date TEXT NOT NULL,
			duration_days INTEGER NOT NULL,
			initial_price_per_kg FLOAT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS bids (
			id SERIAL PRIMARY KEY,
			auction_id INTEGER REFERENCES auctions(id),
			buyer_id INTEGER REFERENCES users(id),
			bid_price_per_kg FLOAT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func IsDuplicateKeyError(err error) bool {
	return errors.Is(err, &pq.Error{Code: "23505"}) // PostgreSQL unique violation code
}
