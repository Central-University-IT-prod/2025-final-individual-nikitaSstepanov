package campaign

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func (c *Campaign) GetClick(ctx ctx.Context, campaignId, clientId string) (*entity.Click, e.Error) {
	query, args := clickQuery(campaignId, clientId)

	row := c.postgres.QueryRow(ctx, query, args...)

	var click entity.Click

	if err := click.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Click wasn`t found.", e.NotFound)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	return &click, nil
}

func (c *Campaign) CreateClick(ctx ctx.Context, click *entity.Click) e.Error {
	query, args := createClickQuery(click)

	tx, err := c.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if err := tx.Commit(ctx); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}
