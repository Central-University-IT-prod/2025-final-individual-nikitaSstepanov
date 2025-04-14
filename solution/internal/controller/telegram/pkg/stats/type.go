package stats

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type TelegramUseCase interface {
	GetState(c ctx.Context, id uint64) (string, e.Error)
	SetState(c ctx.Context, id uint64, state string) e.Error
	GetSession(c ctx.Context, tgId uint64) (string, e.Error)
}

type CampaignUseCase interface {
	AdvertiserBilling(ctx ctx.Context, advertiserId string) (*entity.Billing, e.Error)
	AdvertiserDailyBillWithPagination(ctx ctx.Context, advertiserId string, size, page int) ([]*entity.DailyBilling, e.Error)
	AdvertiserDailyBill(ctx ctx.Context, advertiserId string) ([]*entity.DailyBilling, e.Error)
}
