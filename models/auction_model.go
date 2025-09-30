package models

type Auction struct {
	ID                int     `json:"id"`
	LotID             int     `json:"lot_id"`
	StartDate         string  `json:"start_date"`
	DurationDays      int     `json:"duration_days"`
	InitialPricePerKG float64 `json:"initial_price_per_kg"`
}

type CreateAuctionRequest struct {
	LotID             int     `json:"lot_id"`
	StartDate         string  `json:"start_date"`
	DurationDays      int     `json:"duration_days"`
	InitialPricePerKG float64 `json:"initial_price_per_kg"`
}
