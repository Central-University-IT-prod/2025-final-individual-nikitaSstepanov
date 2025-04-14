package campaign

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type CampaignUseCase interface {
	GetById(ctx ctx.Context, campaign *entity.Campaign) (*entity.Campaign, e.Error)
	Get(ctx ctx.Context, advertiserId string, size int, page int) ([]*entity.Campaign, e.Error)
	UploadImage(ctx ctx.Context, campaign *entity.Campaign, image *entity.Image) (string, e.Error)
	DeleteImage(ctx ctx.Context, campaign *entity.Campaign) e.Error
	GetImage(ctx ctx.Context, campaign *entity.Campaign) (string, e.Error)
	Create(ctx ctx.Context, campaign *entity.Campaign, genText bool, prompt string) e.Error
	Update(ctx ctx.Context, campaign *entity.Campaign, genText bool, prompt string) e.Error
	Delete(ctx ctx.Context, campaign *entity.Campaign) e.Error
}
