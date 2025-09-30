package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"banana-auction/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db                   *sql.DB
	ErrNotFound          = errors.New("record not found")
	ErrDuplicateUsername = errors.New("username already exists")
)

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
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
		return 0, err
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
		return models.User{}, err
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
		return models.User{}, err
	}
	return user, nil
}

func CreateLot(lot models.Lot) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO lots (seller_id, cultivar, planted_country, harvest_date, total_weight_kg)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		lot.SellerID, lot.Cultivar, lot.PlantedCountry, lot.HarvestDate, lot.TotalWeightKG,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetLot(id int) (models.Lot, error) {
	var lot models.Lot
	err := db.QueryRow(`
		SELECT id, seller_id, cultivar, planted_country, harvest_date, total_weight_kg
		FROM lots WHERE id = $1`, id,
	).Scan(&lot.ID, &lot.SellerID, &lot.Cultivar, &lot.PlantedCountry, &lot.HarvestDate, &lot.TotalWeightKG)
	if err == sql.ErrNoRows {
		return models.Lot{}, ErrNotFound
	}
	if err != nil {
		return models.Lot{}, err
	}
	return lot, nil
}

func UpdateLot(lot models.Lot) error {
	_, err := db.Exec(`
		UPDATE lots SET harvest_date = $1
		WHERE id = $2 AND seller_id = $3`,
		lot.HarvestDate, lot.ID, lot.SellerID,
	)
	return err
}

func DeleteLot(id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM bids WHERE auction_id IN (SELECT id FROM auctions WHERE lot_id = $1)`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM auctions WHERE lot_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM lots WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CreateAuction(auction models.Auction) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO auctions (lot_id, start_date, duration_days, initial_price_per_kg)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		auction.LotID, auction.StartDate, auction.DurationDays, auction.InitialPricePerKG,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAuction(id int) (models.Auction, error) {
	var auction models.Auction
	err := db.QueryRow(`
		SELECT id, lot_id, start_date, duration_days, initial_price_per_kg
		FROM auctions WHERE id = $1`, id,
	).Scan(&auction.ID, &auction.LotID, &auction.StartDate, &auction.DurationDays, &auction.InitialPricePerKG)
	if err == sql.ErrNoRows {
		return models.Auction{}, ErrNotFound
	}
	if err != nil {
		return models.Auction{}, err
	}
	return auction, nil
}

func AuctionExistsForLot(lotID int) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM auctions WHERE lot_id = $1`, lotID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateBid(bid models.Bid) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO bids (auction_id, buyer_id, bid_price_per_kg)
		VALUES ($1, $2, $3) RETURNING id`,
		bid.AuctionID, bid.BuyerID, bid.BidPricePerKG,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ListBids(auctionID int) ([]models.Bid, error) {
	rows, err := db.Query(`
		SELECT id, auction_id, buyer_id, bid_price_per_kg
		FROM bids WHERE auction_id = $1`, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.AuctionID, &bid.BuyerID, &bid.BidPricePerKG); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}
