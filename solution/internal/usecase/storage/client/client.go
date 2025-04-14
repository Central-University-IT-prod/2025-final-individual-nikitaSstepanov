package client

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Client struct {
	postgres pg.Client
}

func New(postgres pg.Client) *Client {
	return &Client{
		postgres: postgres,
	}
}

func (c *Client) GetById(ctx ctx.Context, id string) (*entity.Client, e.Error) {
	var client entity.Client

	query, args := idQuery(id)

	row := c.postgres.QueryRow(ctx, query, args...)

	if err := client.Scan(row); err != nil {
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

	return &client, nil
}

func (c *Client) Create(ctx ctx.Context, clients []*entity.Client) e.Error {
	query, args := createQuery(clients)

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

func (c *Client) Update(ctx ctx.Context, client *entity.Client) e.Error {
	query, args := updateQuery(client)

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
