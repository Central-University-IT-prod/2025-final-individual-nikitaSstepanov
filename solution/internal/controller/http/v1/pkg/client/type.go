package client

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type ClientUseCase interface {
	GetById(c ctx.Context, id string) (*entity.Client, e.Error)
	Bulk(c ctx.Context, clients []*entity.Client) e.Error
}
