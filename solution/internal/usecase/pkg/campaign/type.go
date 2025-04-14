package campaign

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type CampaignStorage interface {
	GetById(ctx ctx.Context, id string) (*entity.Campaign, e.Error)
	GetAvailable(ctx ctx.Context, client *entity.Client, day int) ([]*entity.Campaign, e.Error)
	Get(ctx ctx.Context, advertiserId string) ([]*entity.Campaign, e.Error)
	GetWithPagination(ctx ctx.Context, advertiserId string, limit int, offset int) ([]*entity.Campaign, e.Error)
	Create(ctx ctx.Context, campaign *entity.Campaign) e.Error
	Update(ctx ctx.Context, campaign *entity.Campaign) e.Error
	Delete(ctx ctx.Context, campaign *entity.Campaign) e.Error
	ClickStorage
	ImpressionStorage
	BillingStorage
}

type ClickStorage interface {
	GetClick(ctx ctx.Context, campaignId, clientId string) (*entity.Click, e.Error)
	CreateClick(ctx ctx.Context, click *entity.Click) e.Error
}

type ImpressionStorage interface {
	CreateImpression(ctx ctx.Context, impression *entity.Impression) e.Error
	GetImpression(ctx ctx.Context, campaignId, clientId string) (*entity.Impression, e.Error)
}

type BillingStorage interface {
	GetDailyBillingWithPagination(ctx ctx.Context, id string, limit, offset int) ([]*entity.DailyBilling, e.Error)
	GetBilling(ctx ctx.Context, id string) (*entity.Billing, e.Error)
	GetDailyBilling(ctx ctx.Context, id string) ([]*entity.DailyBilling, e.Error)
	UpdateBilling(ctx ctx.Context, billing *entity.Billing, id string) e.Error
	GetBillingByDay(ctx ctx.Context, id string, day int) (*entity.DailyBilling, e.Error)
	CreateDailyBilling(ctx ctx.Context, billing *entity.DailyBilling, id string) e.Error
	UpdateDailyBilling(ctx ctx.Context, billing *entity.DailyBilling, id string) e.Error
}

type TimeUseCase interface {
	Get(c ctx.Context) (int, e.Error)
}

type AdvertiserStorage interface {
	GetById(c ctx.Context, id string) (*entity.Advertiser, e.Error)
}

type ClientStorage interface {
	GetById(ctx ctx.Context, id string) (*entity.Client, e.Error)
}

type ImageStorage interface {
	Delete(c ctx.Context, name string) e.Error
	Download(c ctx.Context, name string) ([]byte, e.Error)
	Get(c ctx.Context, name string) (string, e.Error)
	Upload(c ctx.Context, image *entity.Image) (string, e.Error)
}

type ScoreStorage interface {
	GetClientScores(c ctx.Context, id string) (map[string]int, e.Error)
}

type AiUseCase interface {
	GenText(prompt string) (string, e.Error)
}

type ModerateUseCase interface {
	Moderate(c ctx.Context, text string) e.Error
}
