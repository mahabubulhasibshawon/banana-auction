package lot

import "errors"

type Service interface {
	CreateLot(sellerID int, cultivar, plantedCountry, harvestDate string, totalWeightKG int) (int, error)
	GetLot(id int) (Lot, error)
	UpdateLot(id int, sellerID int, harvestDate string) error
	DeleteLot(id int, sellerID int) error
	ListLots() ([]Lot, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateLot(sellerID int, cultivar, plantedCountry, harvestDate string, totalWeightKG int) (int, error) {
	if totalWeightKG < 1000 {
		return 0, errors.New("minimum weight allowed is 1000 kg")
	}

	l := Lot{
		SellerID:       sellerID,
		Cultivar:       cultivar,
		PlantedCountry: plantedCountry,
		HarvestDate:    harvestDate,
		TotalWeightKG:  totalWeightKG,
	}

	return s.repo.Create(l)
}

func (s *service) GetLot(id int) (Lot, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateLot(id int, sellerID int, harvestDate string) error {
	l, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if l.SellerID != sellerID {
		return errors.New("unauthorized to update this lot")
	}
	l.HarvestDate = harvestDate
	return s.repo.Update(l)
}

func (s *service) DeleteLot(id int, sellerID int) error {
	l, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if l.SellerID != sellerID {
		return errors.New("unauthorized to delete this lot")
	}
	return s.repo.Delete(id)
}

func (s *service) ListLots() ([]Lot, error) {
	return s.repo.List()
}
