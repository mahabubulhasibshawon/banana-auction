package auction

import "errors"

type Service interface {
	CreateAuction(lotID int, startDate string, durationDays int, initialPricePerKG float64) (int, error)
	GetAuction(id int) (Auction, error)
	UpdateAuction(id int, startDate string, durationDays int, initialPricePerKG float64) error
	DeleteAuction(id int) error
	ListAuctions() ([]Auction, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAuction(lotID int, startDate string, durationDays int, initialPricePerKG float64) (int, error) {
	exists, err := s.repo.ExistsForLot(lotID)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("auction already exists for this lot")
	}

	a := Auction{
		LotID:             lotID,
		StartDate:         startDate,
		DurationDays:      durationDays,
		InitialPricePerKG: initialPricePerKG,
	}

	return s.repo.Create(a)
}

func (s *service) GetAuction(id int) (Auction, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateAuction(id int, startDate string, durationDays int, initialPricePerKG float64) error {
	a, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	a.StartDate = startDate
	a.DurationDays = durationDays
	a.InitialPricePerKG = initialPricePerKG
	return s.repo.Update(a)
}

func (s *service) DeleteAuction(id int) error {
	return s.repo.Delete(id)
}

func (s *service) ListAuctions() ([]Auction, error) {
	return s.repo.List()
}
