package new

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type TelegramUseCase interface {
	GetState(c ctx.Context, id uint64) (string, e.Error)
	SetState(c ctx.Context, id uint64, state string) e.Error
	GetSession(c ctx.Context, tgId uint64) (string, e.Error)
	GetNew(c ctx.Context, tgId uint64) (*entity.CampaignData, e.Error)
	SetNew(c ctx.Context, tgId uint64, data *entity.CampaignData) e.Error
	DeleteNew(c ctx.Context, tgId uint64) e.Error
}

type CampaignUseCase interface {
	UploadImage(ctx ctx.Context, campaign *entity.Campaign, image *entity.Image) (string, e.Error)
	Create(ctx ctx.Context, campaign *entity.Campaign, genText bool, prompt string) e.Error
}

type AiUseCase interface {
	GenText(prompt string) (string, e.Error)
}
