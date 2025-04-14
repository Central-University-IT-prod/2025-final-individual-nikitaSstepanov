package client

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type ClientStorage interface {
	GetById(ctx ctx.Context, id string) (*entity.Client, e.Error)
	Create(ctx ctx.Context, client []*entity.Client) e.Error
	Update(ctx ctx.Context, client *entity.Client) e.Error
}
