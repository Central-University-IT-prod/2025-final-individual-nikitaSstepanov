package dto

import "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"

type CreateCampaign struct {
	Title             *string    `json:"ad_title"            validate:"required"`
	Text              *string    `json:"ad_text"             validate:"required"`
	StartDate         *int       `json:"start_date"          validate:"required,gte=0"`
	EndDate           *int       `json:"end_date"            validate:"required,gte=0"`
	ImpressionsLimit  *int       `json:"impressions_limit"   validate:"required,gte=0"`
	ClicksLimit       *int       `json:"clicks_limit"        validate:"required,gte=0"`
	CostPerImpression *float32   `json:"cost_per_impression" validate:"required,gte=0"`
	CostPerClick      *float32   `json:"cost_per_click"      validate:"required,gte=0"`
	GenText           *bool      `json:"gen_text,omitempty"  validate:"omitempty"`
	Prompt            *string    `json:"prompt,omitempty"    validate:"omitempty"`
	Targeting         *Targeting `json:"targeting,omitempty"`
}

type Targeting struct {
	Gender   *types.Gender `json:"gender,omitempty"   validate:"omitempty,oneof=MALE FEMALE ALL"`
	AgeFrom  *int          `json:"age_from,omitempty" validate:"omitempty,gte=0,lte=120"`
	AgeTo    *int          `json:"age_to,omitempty"   validate:"omitempty,gte=0,lte=120"`
	Location *string       `json:"location,omitempty"`
}

type UpdateCampaign struct {
	Title             *string    `json:"ad_title"            validate:"required"`
	Text              *string    `json:"ad_text"             validate:"required"`
	StartDate         *int       `json:"start_date"          validate:"required,gte=0"`
	EndDate           *int       `json:"end_date"            validate:"required,gte=0"`
	ImpressionsLimit  *int       `json:"impressions_limit"   validate:"required,gte=0"`
	ClicksLimit       *int       `json:"clicks_limit"        validate:"required,gte=0"`
	CostPerImpression *float32   `json:"cost_per_impression" validate:"required,gte=0"`
	CostPerClick      *float32   `json:"cost_per_click"      validate:"required,gte=0"`
	GenText           *bool      `json:"gen_text,omitempty"  validate:"omitempty"`
	Prompt            *string    `json:"prompt,omitempty"    validate:"omitempty"`
	Targeting         *Targeting `json:"targeting"           validate:"required"`
}

type ImageUrl struct {
	Url string `json:"image_url"`
}
