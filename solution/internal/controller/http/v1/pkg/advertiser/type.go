package advertiser

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type AdvertiserUseCase interface {
	GetById(c ctx.Context, id string) (*entity.Advertiser, e.Error)
	Bulk(c ctx.Context, advertisers []*entity.Advertiser) e.Error
}
