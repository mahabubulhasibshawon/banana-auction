package bid

type Repository interface {
	Create(b Bid) (int, error)
	GetByID(id int) (Bid, error)
	Update(b Bid) error
	Delete(id int) error
	ListByAuctionID(auctionID int) ([]Bid, error)
}