package converter

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func DtoAdvertiser(advertiser *entity.Advertiser) *dto.Advertiser {
	return &dto.Advertiser{
		Id:   advertiser.Id,
		Name: advertiser.Name,
	}
}

func BulkAdvertiser(body []dto.Advertiser) []*entity.Advertiser {
	result := make([]*entity.Advertiser, 0)

	for _, advertiser := range body {
		toAdd := &entity.Advertiser{
			Id:   advertiser.Id,
			Name: advertiser.Name,
		}

		result = append(result, toAdd)
	}

	return result
}
