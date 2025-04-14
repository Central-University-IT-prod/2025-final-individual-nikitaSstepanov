package advertiser

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type AdvertiserStorage interface {
	GetById(c ctx.Context, id string) (*entity.Advertiser, e.Error)
	Create(ctx ctx.Context, advertisers []*entity.Advertiser) e.Error
	Update(ctx ctx.Context, advertiser *entity.Advertiser) e.Error
}
