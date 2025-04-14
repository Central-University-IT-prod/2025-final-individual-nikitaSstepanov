package dto

type MlScore struct {
	ClientId     string `json:"client_id"     validate:"required,uuid"`
	AdvertiserId string `json:"advertiser_id" validate:"required,uuid"`
	Score        int    `json:"score"         validate:"required"`
}
