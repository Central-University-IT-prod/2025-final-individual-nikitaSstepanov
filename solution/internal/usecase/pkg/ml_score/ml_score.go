package mlscore

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type MlScore struct {
	storage ScoreStorage
}

func New(storage ScoreStorage) *MlScore {
	return &MlScore{
		storage: storage,
	}
}

func (ml *MlScore) Manage(c ctx.Context, score *entity.MlScore) e.Error {
	_, err := ml.storage.GetById(c, score.ClientId, score.AdvertiserId)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	if err != nil && err.GetCode() == e.NotFound {
		return ml.storage.Create(c, score)
	} else {
		return ml.storage.Update(c, score)
	}
}

func (ml *MlScore) GetClientScores(c ctx.Context, id string) (map[string]int, e.Error) {
	scores, err := ml.storage.GetByClient(c, id)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)

	for _, score := range scores {
		result[score.AdvertiserId] = score.Score
	}

	return result, nil
}
