package session

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
}

type AdvertiserUseCase interface {
	GetById(c ctx.Context, id string) (*entity.Advertiser, e.Error)
}
