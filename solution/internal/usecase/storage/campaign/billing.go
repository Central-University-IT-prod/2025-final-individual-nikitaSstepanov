package campaign

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func (c *Campaign) GetBilling(ctx ctx.Context, id string) (*entity.Billing, e.Error) {
	query, args := billQuery(id)

	row := c.postgres.QueryRow(ctx, query, args...)

	var billing entity.Billing

	if err := billing.Scan(row); err != nil {
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

	return &billing, nil
}

func (c *Campaign) GetDailyBilling(ctx ctx.Context, id string) ([]*entity.DailyBilling, e.Error) {
	query, args := dailyBillQuery(id)

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

	bills := make([]*entity.DailyBilling, 0)

	for rows.Next() {
		var bill entity.DailyBilling

		if err := bill.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}

		bills = append(bills, &bill)
	}

	return bills, nil
}

func (c *Campaign) GetDailyBillingWithPagination(ctx ctx.Context, id string, limit, offset int) ([]*entity.DailyBilling, e.Error) {
	query, args := dailyPaginationQuery(id, limit, offset)

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

	bills := make([]*entity.DailyBilling, 0)

	for rows.Next() {
		var bill entity.DailyBilling

		if err := bill.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}

		bills = append(bills, &bill)
	}

	return bills, nil
}

func (c *Campaign) UpdateBilling(ctx ctx.Context, billing *entity.Billing, id string) e.Error {
	query, args := updateCountsQuery(billing, id)

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

func (c *Campaign) GetBillingByDay(ctx ctx.Context, id string, day int) (*entity.DailyBilling, e.Error) {
	query, args := dailyQuery(id, day)

	row := c.postgres.QueryRow(ctx, query, args...)

	var billing entity.DailyBilling

	if err := billing.Scan(row); err != nil {
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

	return &billing, nil
}

func (c *Campaign) CreateDailyBilling(ctx ctx.Context, billing *entity.DailyBilling, id string) e.Error {
	query, args := createDailyQuery(billing, id)

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

func (c *Campaign) UpdateDailyBilling(ctx ctx.Context, billing *entity.DailyBilling, id string) e.Error {
	query, args := updateDailyQuery(billing, id)

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
