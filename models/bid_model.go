package models

type Bid struct {
	ID            int     `json:"id"`
	AuctionID     int     `json:"auction_id"`
	BuyerID       int     `json:"buyer_id"`
	BidPricePerKG float64 `json:"bid_price_per_kg"`
}

type CreateBidRequest struct {
	BidPricePerKG float64 `json:"bid_price_per_kg"`
}
