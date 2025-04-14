package dto

type Advertiser struct {
	Id   string `json:"advertiser_id" validate:"required,uuid"`
	Name string `json:"name"          validate:"required"`
}
