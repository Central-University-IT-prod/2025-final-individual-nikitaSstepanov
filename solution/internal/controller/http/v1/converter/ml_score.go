package converter

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func EntityScore(score dto.MlScore) *entity.MlScore {
	return &entity.MlScore{
		ClientId:     score.ClientId,
		AdvertiserId: score.AdvertiserId,
		Score:        score.Score,
	}
}
