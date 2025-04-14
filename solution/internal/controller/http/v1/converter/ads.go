package converter

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func DtoAd(campaign *entity.Campaign) map[string]interface{} {
	data := make(map[string]interface{})

	data["ad_id"] = campaign.Id
	data["ad_title"] = campaign.Title
	data["ad_text"] = campaign.Text
	data["advertiser_id"] = campaign.AdvertiserId

	if campaign.Image != "" {
		data["image_url"] = campaign.Image
	}

	return data
}
