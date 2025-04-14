package time

import (
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type Time struct {
	redis rs.Client
}

func New(c ctx.Context, redis rs.Client) *Time {
	storage := &Time{
		redis: redis,
	}

	if err := storage.Set(c, 0); err != nil {
		log := c.Logger()

		log.Error("Can`t set time.", err.SlErr())
		panic("App start error.")
	}

	return storage
}

func (t *Time) Get(c ctx.Context) (int, e.Error) {
	var result int

	err := t.redis.Get(c, redisKey()).Scan(&result)
	if err != nil {
		if err == rs.Nil {
			return 0, notFoundErr.
				WithErr(err).
				WithCtx(c)
		} else {
			return 0, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return result, nil
}

func (t *Time) Set(c ctx.Context, day int) e.Error {
	err := t.redis.Set(c, redisKey(), day, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}
