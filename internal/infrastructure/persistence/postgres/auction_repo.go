package postgres

import (
	"database/sql"
	"errors"

	"banana-auction/internal/domain/auction"
)

type AuctionRepo struct {
	db *sql.DB
}

func NewAuctionRepo(db *sql.DB) *AuctionRepo {
	return &AuctionRepo{db: db}
}

func (r *AuctionRepo) Create(a auction.Auction) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO auctions (lot_id, start_date, duration_days, initial_price_per_kg)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		a.LotID, a.StartDate, a.DurationDays, a.InitialPricePerKG,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuctionRepo) GetByID(id int) (auction.Auction, error) {
	var a auction.Auction
	err := r.db.QueryRow(`
		SELECT id, lot_id, start_date, duration_days, initial_price_per_kg
		FROM auctions WHERE id = $1`, id,
	).Scan(&a.ID, &a.LotID, &a.StartDate, &a.DurationDays, &a.InitialPricePerKG)
	if err == sql.ErrNoRows {
		return auction.Auction{}, errors.New("auction not found")
	}
	if err != nil {
		return auction.Auction{}, err
	}
	return a, nil
}

func (r *AuctionRepo) Update(a auction.Auction) error {
	_, err := r.db.Exec(`
		UPDATE auctions SET start_date = $1, duration_days = $2, initial_price_per_kg = $3
		WHERE id = $4`,
		a.StartDate, a.DurationDays, a.InitialPricePerKG, a.ID,
	)
	return err
}

func (r *AuctionRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM auctions WHERE id = $1`, id)
	return err
}

func (r *AuctionRepo) List() ([]auction.Auction, error) {
	rows, err := r.db.Query(`
		SELECT id, lot_id, start_date, duration_days, initial_price_per_kg
		FROM auctions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []auction.Auction
	for rows.Next() {
		var a auction.Auction
		if err := rows.Scan(&a.ID, &a.LotID, &a.StartDate, &a.DurationDays, &a.InitialPricePerKG); err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}
	return auctions, nil
}

func (r *AuctionRepo) ExistsForLot(lotID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM auctions WHERE lot_id = $1`, lotID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
