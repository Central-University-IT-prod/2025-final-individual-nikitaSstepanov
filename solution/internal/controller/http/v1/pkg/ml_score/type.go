package mlscore

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type ScoreUseCase interface {
	Manage(c ctx.Context, score *entity.MlScore) e.Error
}
