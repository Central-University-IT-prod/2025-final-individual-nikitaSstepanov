package mlscore

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type ScoreStorage interface {
	GetByClient(c ctx.Context, clientId string) ([]*entity.MlScore, e.Error)
	GetById(c ctx.Context, clientId string, advertiserId string) (*entity.MlScore, e.Error)
	Create(ctx ctx.Context, score *entity.MlScore) e.Error
	Update(ctx ctx.Context, score *entity.MlScore) e.Error
}
