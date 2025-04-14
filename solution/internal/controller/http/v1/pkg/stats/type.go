package stats

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type BillingUseCase interface {
	CampaignBilling(ctx ctx.Context, id string) (*entity.Billing, e.Error)
	AdvertiserBilling(ctx ctx.Context, advertiserId string) (*entity.Billing, e.Error)
	CampaignDailyBilling(ctx ctx.Context, id string) ([]*entity.DailyBilling, e.Error)
	AdvertiserDailyBill(ctx ctx.Context, advertiserId string) ([]*entity.DailyBilling, e.Error)
}
