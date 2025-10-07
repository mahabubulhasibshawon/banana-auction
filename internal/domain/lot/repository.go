package lot

type Repository interface {
	Create(l Lot) (int, error)
	GetByID(id int) (Lot, error)
	Update(l Lot) error
	Delete(id int) error
	List() ([]Lot, error)
}