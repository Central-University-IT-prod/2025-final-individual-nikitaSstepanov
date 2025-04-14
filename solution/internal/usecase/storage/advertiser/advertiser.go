package advertiser

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Advertiser struct {
	postgres pg.Client
}

func New(postgres pg.Client) *Advertiser {
	return &Advertiser{
		postgres: postgres,
	}
}

func (a *Advertiser) GetById(c ctx.Context, id string) (*entity.Advertiser, e.Error) {
	query, args := idQuery(id)

	row := a.postgres.QueryRow(c, query, args...)

	var advertiser entity.Advertiser

	if err := advertiser.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &advertiser, nil
}

func (a *Advertiser) Create(ctx ctx.Context, advertisers []*entity.Advertiser) e.Error {
	query, args := createQuery(advertisers)

	tx, err := a.postgres.Begin(ctx)
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

func (a *Advertiser) Update(ctx ctx.Context, advertiser *entity.Advertiser) e.Error {
	query, args := updateQuery(advertiser)

	tx, err := a.postgres.Begin(ctx)
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
