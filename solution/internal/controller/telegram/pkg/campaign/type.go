package campaign

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type TelegramUseCase interface {
	GetState(c ctx.Context, id uint64) (string, e.Error)
	SetState(c ctx.Context, id uint64, state string) e.Error
	GetSession(c ctx.Context, tgId uint64) (string, e.Error)
	SetSession(c ctx.Context, tgId uint64, id string) e.Error
	DeleteSession(c ctx.Context, tgId uint64) e.Error
	GetNew(c ctx.Context, tgId uint64) (*entity.CampaignData, e.Error)
	SetNew(c ctx.Context, tgId uint64, data *entity.CampaignData) e.Error
	DeleteNew(c ctx.Context, tgId uint64) e.Error
	GetCampaignsCount(c ctx.Context, tgId uint64) (int, e.Error)
	GetCampaignId(c ctx.Context, tgId uint64) (string, e.Error)
	SetCampaignId(c ctx.Context, tgId uint64, id string) e.Error
	DeleteCampaignId(c ctx.Context, tgId uint64) e.Error
}

type AiUseCase interface {
	GenText(prompt string) (string, e.Error)
}

type CampaignUseCase interface {
	CampaignDailyBillingWithPagination(ctx ctx.Context, id string, size, page int) ([]*entity.DailyBilling, e.Error)
	CampaignDailyBilling(ctx ctx.Context, id string) ([]*entity.DailyBilling, e.Error)
	CampaignBilling(ctx ctx.Context, id string) (*entity.Billing, e.Error)
	GetById(ctx ctx.Context, campaign *entity.Campaign) (*entity.Campaign, e.Error)
	Get(ctx ctx.Context, advertiserId string, size int, page int) ([]*entity.Campaign, e.Error)
	DownloadImage(ctx ctx.Context, campaign *entity.Campaign) ([]byte, e.Error)
	Delete(ctx ctx.Context, campaign *entity.Campaign) e.Error
	Update(ctx ctx.Context, campaign *entity.Campaign, genText bool, prompt string) e.Error
	DeleteImage(ctx ctx.Context, campaign *entity.Campaign) e.Error
	UploadImage(ctx ctx.Context, campaign *entity.Campaign, image *entity.Image) (string, e.Error)
}

type TimeUseCase interface {
	Get(c ctx.Context) (int, e.Error)
}
