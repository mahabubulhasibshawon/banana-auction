package postgres

import (
	"database/sql"
	"errors"

	"banana-auction/internal/domain/bid"
)

type BidRepo struct {
	db *sql.DB
}

func NewBidRepo(db *sql.DB) *BidRepo {
	return &BidRepo{db: db}
}

func (r *BidRepo) Create(b bid.Bid) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO bids (auction_id, buyer_id, bid_price_per_kg)
		VALUES ($1, $2, $3) RETURNING id`,
		b.AuctionID, b.BuyerID, b.BidPricePerKG,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *BidRepo) GetByID(id int) (bid.Bid, error) {
	var b bid.Bid
	err := r.db.QueryRow(`
		SELECT id, auction_id, buyer_id, bid_price_per_kg
		FROM bids WHERE id = $1`, id,
	).Scan(&b.ID, &b.AuctionID, &b.BuyerID, &b.BidPricePerKG)
	if err == sql.ErrNoRows {
		return bid.Bid{}, errors.New("bid not found")
	}
	if err != nil {
		return bid.Bid{}, err
	}
	return b, nil
}

func (r *BidRepo) Update(b bid.Bid) error {
	_, err := r.db.Exec(`
		UPDATE bids SET bid_price_per_kg = $1
		WHERE id = $2`,
		b.BidPricePerKG, b.ID,
	)
	return err
}

func (r *BidRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM bids WHERE id = $1`, id)
	return err
}

func (r *BidRepo) ListByAuctionID(auctionID int) ([]bid.Bid, error) {
	rows, err := r.db.Query(`
		SELECT id, auction_id, buyer_id, bid_price_per_kg
		FROM bids WHERE auction_id = $1`, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []bid.Bid
	for rows.Next() {
		var b bid.Bid
		if err := rows.Scan(&b.ID, &b.AuctionID, &b.BuyerID, &b.BidPricePerKG); err != nil {
			return nil, err
		}
		bids = append(bids, b)
	}
	return bids, nil
}
