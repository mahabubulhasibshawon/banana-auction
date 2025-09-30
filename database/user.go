package database

import (
	"database/sql"
	"errors"
	"fmt"

	"banana-auction/config"
	"banana-auction/models"

	_ "github.com/lib/pq"
)

var (
	db                   *sql.DB
	ErrNotFound          = errors.New("record not found")
	ErrDuplicateUsername = errors.New("username already exists")
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

func CreateUser(user models.User) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO users (username, password_hash, name, role)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Username, user.PasswordHash, user.Name, user.Role,
	).Scan(&id)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return 0, ErrDuplicateUsername
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, username, password_hash, name, role
		FROM users WHERE username = $1`, username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

func GetUser(id int) (models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, username, password_hash, name, role
		FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
