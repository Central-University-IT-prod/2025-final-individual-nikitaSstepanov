package dto

type Ad struct {
	Id           string `json:"ad_id"`
	Title        string `json:"ad_title"`
	Text         string `json:"ad_text"`
	AdvertiserId string `json:"advertiser_id"`
}

type ClientId struct {
	Id string `json:"client_id" validate:"required,uuid"`
}
