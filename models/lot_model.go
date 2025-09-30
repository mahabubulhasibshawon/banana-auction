package models

type Lot struct {
	ID             int    `json:"id"`
	SellerID       int    `json:"seller_id"`
	Cultivar       string `json:"cultivar"`
	PlantedCountry string `json:"planted_country"`
	HarvestDate    string `json:"harvest_date"`
	TotalWeightKG  int    `json:"total_weight_kg"`
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
