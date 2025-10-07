package auction

type Auction struct {
	ID                int
	LotID             int
	StartDate         string
	DurationDays      int
	InitialPricePerKG float64
}