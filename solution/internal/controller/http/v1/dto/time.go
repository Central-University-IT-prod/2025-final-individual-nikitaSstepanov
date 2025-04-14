package dto

type AdvanceTime struct {
	Day int `json:"current_date" validate:"required,gte=0"`
}
