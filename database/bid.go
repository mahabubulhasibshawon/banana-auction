package database

import (
	"banana-auction/models"
	"database/sql"
	"fmt"
)

func CreateBid(bid models.Bid) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO bids (auction_id, buyer_id, bid_price_per_kg)
		VALUES ($1, $2, $3) RETURNING id, bid_price_per_kg`,
		bid.AuctionID, bid.BuyerID, bid.BidPricePerKG,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create bid: %w", err)
	}
	return id, nil
}

func GetBid(id int) (models.Bid, error) {
	var bid models.Bid
	err := db.QueryRow(`
		SELECT id, auction_id, buyer_id, bid_price_per_kg
		FROM bids WHERE id = $1`, id,
	).Scan(&bid.ID, &bid.AuctionID, &bid.BuyerID, &bid.BidPricePerKG)
	if err == sql.ErrNoRows {
		return models.Bid{}, ErrNotFound
	}
	if err != nil {
		return models.Bid{}, fmt.Errorf("failed to get bid: %w", err)
	}
	return bid, nil
}

func UpdateBid(bid models.Bid) error {
	_, err := db.Exec(`
		UPDATE bids SET bid_price_per_kg = $1
		WHERE id = $2`,
		bid.BidPricePerKG, bid.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update bid: %w", err)
	}
	return nil
}

func DeleteBid(id int) error {
	_, err := db.Exec(`DELETE FROM bids WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete bid: %w", err)
	}
	return nil
}

func ListBids(auctionID int) ([]models.Bid, error) {
	rows, err := db.Query(`
		SELECT id, auction_id, buyer_id, bid_price_per_kg
		FROM bids WHERE auction_id = $1`, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bids: %w", err)
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.AuctionID, &bid.BuyerID, &bid.BidPricePerKG); err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, bid)
	}
	return bids, nil
}
