package entity

import (
	"encoding/json"

	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"
)

type CampaignData struct {
	AdvertiserId      string       `redis:"advertiser_id"`
	Title             string       `redis:"title"`
	Text              string       `redis:"text"`
	StartDate         int          `redis:"start_date"`
	EndDate           int          `redis:"end_date"`
	ImpressionsLimit  int          `reids:"impressions_limit"`
	ClicksLimit       int          `redis:"clicks_limit"`
	CostPerImpression float32      `redis:"cost_per_impression"`
	CostPerClick      float32      `redis:"cost_per_click"`
	Gender            *types.Gender `redis:"gender"`
	AgeFrom           *int          `redis:"age_from"`
	AgeTo             *int          `redis:"age_to"`
	Location          *string       `redis:"location"`
}

func (c *CampaignData) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CampaignData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
