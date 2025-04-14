package converter

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func DtoCampaign(campaign *entity.Campaign) map[string]interface{} {
	data := make(map[string]interface{})

	data["campaign_id"] = campaign.Id
	data["advertiser_id"] = campaign.AdvertiserId
	data["impressions_limit"] = campaign.Billing.ImpressionsLimit
	data["clicks_limit"] = campaign.Billing.ClicksLimit
	data["cost_per_impression"] = campaign.Billing.CostPerImpression
	data["cost_per_click"] = campaign.Billing.CostPerClick
	data["ad_title"] = campaign.Title
	data["ad_text"] = campaign.Text
	data["start_date"] = campaign.StartDate
	data["end_date"] = campaign.EndDate
	data["targeting"] = make(map[string]interface{})

	if campaign.Targeting != nil {
		if campaign.Targeting.AgeFrom != nil {
			data["targeting"].(map[string]interface{})["age_from"] = campaign.Targeting.AgeFrom
		}

		if campaign.Targeting.AgeTo != nil {
			data["targeting"].(map[string]interface{})["age_to"] = campaign.Targeting.AgeTo
		}

		if campaign.Targeting.Gender != nil {
			data["targeting"].(map[string]interface{})["gender"] = campaign.Targeting.Gender
		}

		if campaign.Targeting.Location != nil {
			data["targeting"].(map[string]interface{})["location"] = campaign.Targeting.Location
		}
	}

	return data
}

func CreateCampaign(body dto.CreateCampaign, advertiserId string) *entity.Campaign {
	campaign := &entity.Campaign{
		AdvertiserId: advertiserId,
		Title:        *body.Title,
		Text:         *body.Text,
		StartDate:    *body.StartDate,
		EndDate:      *body.EndDate,
	}

	billing := &entity.Billing{
		ImpressionsLimit:  *body.ImpressionsLimit,
		ClicksLimit:       *body.ClicksLimit,
		CostPerImpression: *body.CostPerImpression,
		CostPerClick:      *body.CostPerClick,
	}

	campaign.Billing = billing

	if body.Targeting != nil {
		targeting := &entity.Targeting{
			Gender:   body.Targeting.Gender,
			AgeFrom:  body.Targeting.AgeFrom,
			AgeTo:    body.Targeting.AgeTo,
			Location: body.Targeting.Location,
		}

		campaign.Targeting = targeting
	}

	return campaign
}

func UpdateCampaign(body dto.UpdateCampaign) *entity.Campaign {
	campaign := &entity.Campaign{
		Title:     *body.Title,
		Text:      *body.Text,
		StartDate: *body.StartDate,
		EndDate:   *body.EndDate,
	}

	billing := &entity.Billing{
		ImpressionsLimit:  *body.ImpressionsLimit,
		ClicksLimit:       *body.ClicksLimit,
		CostPerImpression: *body.CostPerImpression,
		CostPerClick:      *body.CostPerClick,
	}

	campaign.Billing = billing

	targeting := &entity.Targeting{
		Gender:   body.Targeting.Gender,
		AgeFrom:  body.Targeting.AgeFrom,
		AgeTo:    body.Targeting.AgeTo,
		Location: body.Targeting.Location,
	}

	campaign.Targeting = targeting

	return campaign
}
