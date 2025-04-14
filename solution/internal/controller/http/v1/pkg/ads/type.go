package ads

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type AdsUseCase interface {
	GetAd(ctx ctx.Context, clientId string) (*entity.Campaign, e.Error)
	Click(ctx ctx.Context, click *entity.Click) e.Error
}
