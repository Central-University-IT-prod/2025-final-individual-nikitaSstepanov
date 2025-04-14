package dto

import "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"

type Client struct {
	Id       string       `json:"client_id" validate:"required,uuid"`
	Login    string       `json:"login"     validate:"required"`
	Age      int          `json:"age"       validate:"required,gte=1,lte=120"`
	Location string       `json:"location"  validate:"required"`
	Gender   types.Gender `json:"gender"    validate:"required,oneof=MALE FEMALE"`
}
