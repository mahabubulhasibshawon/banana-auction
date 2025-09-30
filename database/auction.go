package database

import (
	"banana-auction/models"
	"database/sql"
	"fmt"
)

func CreateAuction(auction models.Auction) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO auctions (lot_id, start_date, duration_days, initial_price_per_kg)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		auction.LotID, auction.StartDate, auction.DurationDays, auction.InitialPricePerKG,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create auction: %w", err)
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
		return models.Auction{}, fmt.Errorf("failed to get auction: %w", err)
	}
	return auction, nil
}

func UpdateAuction(auction models.Auction) error {
	_, err := db.Exec(`
		UPDATE auctions SET start_date = $1, duration_days = $2, initial_price_per_kg = $3
		WHERE id = $4`,
		auction.StartDate, auction.DurationDays, auction.InitialPricePerKG, auction.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update auction: %w", err)
	}
	return nil
}

func DeleteAuction(id int) error {
	_, err := db.Exec(`DELETE FROM auctions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete auction: %w", err)
	}
	return nil
}

func AuctionExistsForLot(lotID int) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM auctions WHERE lot_id = $1`, lotID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check auction existence: %w", err)
	}
	return count > 0, nil
}

func ListAuctions() ([]models.Auction, error) {
	rows, err := db.Query(`
		SELECT id, lot_id, start_date, duration_days, initial_price_per_kg
		FROM auctions`)
	if err != nil {
		return nil, fmt.Errorf("failed to list auctions: %w", err)
	}
	defer rows.Close()

	var auctions []models.Auction
	for rows.Next() {
		var auction models.Auction
		if err := rows.Scan(&auction.ID, &auction.LotID, &auction.StartDate, &auction.DurationDays, &auction.InitialPricePerKG); err != nil {
			return nil, fmt.Errorf("failed to scan auction: %w", err)
		}
		auctions = append(auctions, auction)
	}
	return auctions, nil
}
