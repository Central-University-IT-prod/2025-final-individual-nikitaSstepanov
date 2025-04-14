package mlscore

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type MlScore struct {
	postgres pg.Client
}

func New(postgres pg.Client) *MlScore {
	return &MlScore{
		postgres: postgres,
	}
}

func (ml *MlScore) GetByClient(c ctx.Context, clientId string) ([]*entity.MlScore, e.Error) {
	query, args := clientQuery(clientId)

	rows, err := ml.postgres.Query(c, query, args...)
	if err != nil {
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

	var result []*entity.MlScore

	for rows.Next() {
		var score entity.MlScore

		if err := score.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		result = append(result, &score)
	}

	return result, nil
}

func (ml *MlScore) GetById(c ctx.Context, clientId string, advertiserId string) (*entity.MlScore, e.Error) {
	query, args := idQuery(clientId, advertiserId)

	row := ml.postgres.QueryRow(c, query, args...)

	var score entity.MlScore

	if err := score.Scan(row); err != nil {
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

	return &score, nil
}

func (ml *MlScore) Create(ctx ctx.Context, score *entity.MlScore) e.Error {
	query, args := createQuery(score)

	tx, err := ml.postgres.Begin(ctx)
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

func (ml *MlScore) Update(ctx ctx.Context, score *entity.MlScore) e.Error {
	query, args := updateQuery(score)

	tx, err := ml.postgres.Begin(ctx)
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
