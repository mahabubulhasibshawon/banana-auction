package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Name         string `json:"name"`
	Role         string `json:"role"` // "seller" or "buyer"
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"` // "seller" or "buyer"
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Lot struct {
	ID             int    `json:"id"`
	SellerID       int    `json:"seller_id"`
	Cultivar       string `json:"cultivar"`
	PlantedCountry string `json:"planted_country"`
	HarvestDate    string `json:"harvest_date"`
	TotalWeightKG  int    `json:"total_weight_kg"`
}

type Auction struct {
	ID                int     `json:"id"`
	LotID             int     `json:"lot_id"`
	StartDate         string  `json:"start_date"`
	DurationDays      int     `json:"duration_days"`
	InitialPricePerKG float64 `json:"initial_price_per_kg"`
}

type Bid struct {
	ID            int     `json:"id"`
	AuctionID     int     `json:"auction_id"`
	BuyerID       int     `json:"buyer_id"`
	BidPricePerKG float64 `json:"bid_price_per_kg"`
}

type CreateLotRequest struct {
	Cultivar       string `json:"cultivar"`
	PlantedCountry string `json:"planted_country"`
	HarvestDate    string `json:"harvest_date"`
	TotalWeightKG  int    `json:"total_weight_kg"`
}

type UpdateLotRequest struct {
	HarvestDate string `json:"harvest_date"`
}

type CreateAuctionRequest struct {
	LotID             int     `json:"lot_id"`
	StartDate         string  `json:"start_date"`
	DurationDays      int     `json:"duration_days"`
	InitialPricePerKG float64 `json:"initial_price_per_kg"`
}

type CreateBidRequest struct {
	BidPricePerKG float64 `json:"bid_price_per_kg"`
}

type IDResponse struct {
	ID int `json:"id"`
}