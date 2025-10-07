package bid

type Service interface {
	PlaceBid(auctionID, buyerID int, bidPricePerKG float64) (int, error)
	GetBid(id int) (Bid, error)
	UpdateBid(id int, bidPricePerKG float64) error
	DeleteBid(id int) error
	ListBids(auctionID int) ([]Bid, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) PlaceBid(auctionID, buyerID int, bidPricePerKG float64) (int, error) {
	b := Bid{
		AuctionID:     auctionID,
		BuyerID:       buyerID,
		BidPricePerKG: bidPricePerKG,
	}

	return s.repo.Create(b)
}

func (s *service) GetBid(id int) (Bid, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateBid(id int, bidPricePerKG float64) error {
	b, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	b.BidPricePerKG = bidPricePerKG
	return s.repo.Update(b)
}

func (s *service) DeleteBid(id int) error {
	return s.repo.Delete(id)
}

func (s *service) ListBids(auctionID int) ([]Bid, error) {
	return s.repo.ListByAuctionID(auctionID)
}