package database

import (
	"banana-auction/models"
	"database/sql"
	"fmt"
)

func CreateLot(lot models.Lot) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO lots (seller_id, cultivar, planted_country, harvest_date, total_weight_kg)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		lot.SellerID, lot.Cultivar, lot.PlantedCountry, lot.HarvestDate, lot.TotalWeightKG,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create lot: %w", err)
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
		return models.Lot{}, fmt.Errorf("failed to get lot: %w", err)
	}
	return lot, nil
}

func UpdateLot(lot models.Lot) error {
	_, err := db.Exec(`
		UPDATE lots SET harvest_date = $1
		WHERE id = $2 AND seller_id = $3`,
		lot.HarvestDate, lot.ID, lot.SellerID,
	)
	if err != nil {
		return fmt.Errorf("failed to update lot: %w", err)
	}
	return nil
}

func DeleteLot(id int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM bids WHERE auction_id IN (SELECT id FROM auctions WHERE lot_id = $1)`, id)
	if err != nil {
		return fmt.Errorf("failed to delete bids: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM auctions WHERE lot_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete auctions: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM lots WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete lot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func ListLots() ([]models.Lot, error) {
	rows, err := db.Query(`
		SELECT id, seller_id, cultivar, planted_country, harvest_date, total_weight_kg
		FROM lots`)
	if err != nil {
		return nil, fmt.Errorf("failed to list lots: %w", err)
	}
	defer rows.Close()

	var lots []models.Lot
	for rows.Next() {
		var lot models.Lot
		if err := rows.Scan(&lot.ID, &lot.SellerID, &lot.Cultivar, &lot.PlantedCountry, &lot.HarvestDate, &lot.TotalWeightKG); err != nil {
			return nil, fmt.Errorf("failed to scan lot: %w", err)
		}
		lots = append(lots, lot)
	}
	return lots, nil
}
