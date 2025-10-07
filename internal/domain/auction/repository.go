package auction

type Repository interface {
	Create(a Auction) (int, error)
	GetByID(id int) (Auction, error)
	Update(a Auction) error
	Delete(id int) error
	List() ([]Auction, error)
	ExistsForLot(lotID int) (bool, error)
}