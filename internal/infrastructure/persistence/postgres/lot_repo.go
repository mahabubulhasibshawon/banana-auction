package postgres

import (
	"database/sql"
	"errors"

	"banana-auction/internal/domain/lot"
)

type LotRepo struct {
	db *sql.DB
}

func NewLotRepo(db *sql.DB) *LotRepo {
	return &LotRepo{db: db}
}

func (r *LotRepo) Create(l lot.Lot) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO lots (seller_id, cultivar, planted_country, harvest_date, total_weight_kg)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		l.SellerID, l.Cultivar, l.PlantedCountry, l.HarvestDate, l.TotalWeightKG,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LotRepo) GetByID(id int) (lot.Lot, error) {
	var l lot.Lot
	err := r.db.QueryRow(`
		SELECT id, seller_id, cultivar, planted_country, harvest_date, total_weight_kg
		FROM lots WHERE id = $1`, id,
	).Scan(&l.ID, &l.SellerID, &l.Cultivar, &l.PlantedCountry, &l.HarvestDate, &l.TotalWeightKG)
	if err == sql.ErrNoRows {
		return lot.Lot{}, errors.New("lot not found")
	}
	if err != nil {
		return lot.Lot{}, err
	}
	return l, nil
}

func (r *LotRepo) Update(l lot.Lot) error {
	_, err := r.db.Exec(`
		UPDATE lots SET harvest_date = $1
		WHERE id = $2 AND seller_id = $3`,
		l.HarvestDate, l.ID, l.SellerID,
	)
	return err
}

func (r *LotRepo) Delete(id int) error {
	tx, err := r.db.Begin()
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

func (r *LotRepo) List() ([]lot.Lot, error) {
	rows, err := r.db.Query(`
		SELECT id, seller_id, cultivar, planted_country, harvest_date, total_weight_kg
		FROM lots`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lots []lot.Lot
	for rows.Next() {
		var l lot.Lot
		if err := rows.Scan(&l.ID, &l.SellerID, &l.Cultivar, &l.PlantedCountry, &l.HarvestDate, &l.TotalWeightKG); err != nil {
			return nil, err
		}
		lots = append(lots, l)
	}
	return lots, nil
}