package telegram

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type TelegramStorage interface {
	GetState(c ctx.Context, id uint64) (string, e.Error)
	SetState(c ctx.Context, id uint64, state string) e.Error
	GetSession(c ctx.Context, tgId uint64) (string, e.Error)
	SetSession(c ctx.Context, tgId uint64, id string) e.Error
	DeleteSession(c ctx.Context, tgId uint64) e.Error
	GetNew(c ctx.Context, tgId uint64) (*entity.CampaignData, e.Error)
	SetNew(c ctx.Context, tgId uint64, data *entity.CampaignData) e.Error
	DeleteNew(c ctx.Context, tgId uint64) e.Error
	GetCampaignId(c ctx.Context, tgId uint64) (string, e.Error)
	SetCampaignId(c ctx.Context, tgId uint64, id string) e.Error
	DeleteCampaignId(c ctx.Context, tgId uint64) e.Error
}

type CampaignStorage interface {
	Get(ctx ctx.Context, advertiserId string) ([]*entity.Campaign, e.Error)
}
