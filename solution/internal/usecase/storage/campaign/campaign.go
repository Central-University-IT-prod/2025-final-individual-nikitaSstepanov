package campaign

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Campaign struct {
	postgres pg.Client
}

func New(postgres pg.Client) *Campaign {
	return &Campaign{
		postgres: postgres,
	}
}

func (c *Campaign) GetById(ctx ctx.Context, id string) (*entity.Campaign, e.Error) {
	var campaign entity.Campaign

	query, args := idQuery(id)

	row := c.postgres.QueryRow(ctx, query, args...)

	if err := campaign.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	return &campaign, nil
}

func (c *Campaign) GetAvailable(ctx ctx.Context, client *entity.Client, day int) ([]*entity.Campaign, e.Error) {
	query, args := availableQuery(client, day)

	rows, err := c.postgres.Query(ctx, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	campaigns := make([]*entity.Campaign, 0)

	for rows.Next() {
		var campaign entity.Campaign

		if err := campaign.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}

		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (c *Campaign) Get(ctx ctx.Context, advertiserId string) ([]*entity.Campaign, e.Error) {
	query, args := getQuery(advertiserId)

	rows, err := c.postgres.Query(ctx, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	campaigns := make([]*entity.Campaign, 0)

	for rows.Next() {
		var campaign entity.Campaign

		if err := campaign.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}

		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (c *Campaign) GetWithPagination(ctx ctx.Context, advertiserId string, limit, offset int) ([]*entity.Campaign, e.Error) {
	query, args := paginationQuery(advertiserId, limit, offset)

	rows, err := c.postgres.Query(ctx, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	campaigns := make([]*entity.Campaign, 0)

	for rows.Next() {
		var campaign entity.Campaign

		if err := campaign.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}

		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (c *Campaign) Create(ctx ctx.Context, campaign *entity.Campaign) e.Error {
	query, args := createQuery(campaign)

	tx, err := c.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, query, args...)

	if err := row.Scan(&campaign.Id); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	query, args = createTargetingQuery(campaign.Targeting, campaign.Id)

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if campaign.Billing == nil {
		return badBillingErr
	}

	query, args = createBillingQuery(campaign.Billing, campaign.Id)

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

func (c *Campaign) Update(ctx ctx.Context, campaign *entity.Campaign) e.Error {
	query, args := updateQuery(campaign)

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

	query, args = updateBillingQuery(campaign.Billing, campaign.Id)

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	query, args = updateTargetingQuery(campaign.Targeting, campaign.Id)

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

func (c *Campaign) Delete(ctx ctx.Context, campaign *entity.Campaign) e.Error {
	query, args := deleteQuery(campaign)

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
