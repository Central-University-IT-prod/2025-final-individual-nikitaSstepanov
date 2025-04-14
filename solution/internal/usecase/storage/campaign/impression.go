package campaign

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func (c *Campaign) GetImpression(ctx ctx.Context, campaignId, clientId string) (*entity.Impression, e.Error) {
	query, args := impressionQuery(campaignId, clientId)

	row := c.postgres.QueryRow(ctx, query, args...)

	var impression entity.Impression

	if err := impression.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Impression wasn`t found.", e.NotFound)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	return &impression, nil
}

func (c *Campaign) CreateImpression(ctx ctx.Context, impression *entity.Impression) e.Error {
	query, args := createImpressionQuery(impression)

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
