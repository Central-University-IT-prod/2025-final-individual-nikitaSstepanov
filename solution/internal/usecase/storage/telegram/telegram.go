package telegram

import (
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Telegram struct {
	redis rs.Client
}

func New(redis rs.Client) *Telegram {
	return &Telegram{
		redis: redis,
	}
}

func (t *Telegram) GetState(c ctx.Context, id uint64) (string, e.Error) {
	var result string

	err := t.redis.Get(c, stateQuery(id)).Scan(&result)
	if err != nil {
		if err == rs.Nil {
			return "", e.New("", e.NotFound)
		} else {
			return "", e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return result, nil
}

func (t *Telegram) SetState(c ctx.Context, id uint64, state string) e.Error {
	err := t.redis.Set(c, stateQuery(id), state, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (t *Telegram) GetSession(c ctx.Context, tgId uint64) (string, e.Error) {
	var result string

	err := t.redis.Get(c, sessionQuery(tgId)).Scan(&result)
	if err != nil {
		if err == rs.Nil {
			return "", e.New("", e.NotFound)
		} else {
			return "", e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return result, nil
}

func (t *Telegram) SetSession(c ctx.Context, tgId uint64, id string) e.Error {
	err := t.redis.Set(c, sessionQuery(tgId), id, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (t *Telegram) DeleteSession(c ctx.Context, tgId uint64) e.Error {
	err := t.redis.Del(c, sessionQuery(tgId)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	_, getErr := t.GetState(c, tgId)
	if getErr != nil && getErr.GetCode() != e.NotFound {
		return getErr
	}

	if getErr != nil {
		return nil
	}

	err = t.redis.Del(c, stateQuery(tgId)).Err()
	if err != nil {
		return e.InternalErr.WithErr(err)
	}

	return nil
}

func (t *Telegram) GetNew(c ctx.Context, tgId uint64) (*entity.CampaignData, e.Error) {
	var result entity.CampaignData

	err := t.redis.Get(c, newQuery(tgId)).Scan(&result)
	if err != nil {
		if err == rs.Nil {
			return nil, e.New("", e.NotFound)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &result, nil
}

func (t *Telegram) SetNew(c ctx.Context, tgId uint64, data *entity.CampaignData) e.Error {
	err := t.redis.Set(c, newQuery(tgId), data, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (t *Telegram) DeleteNew(c ctx.Context, tgId uint64) e.Error {
	err := t.redis.Del(c, newQuery(tgId)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (t *Telegram) GetCampaignId(c ctx.Context, tgId uint64) (string, e.Error) {
	var result string

	err := t.redis.Get(c, campaignQuery(tgId)).Scan(&result)
	if err != nil {
		if err == rs.Nil {
			return "", e.New("", e.NotFound)
		} else {
			return "", e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return result, nil
}

func (t *Telegram) SetCampaignId(c ctx.Context, tgId uint64, id string) e.Error {
	err := t.redis.Set(c, campaignQuery(tgId), id, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (t *Telegram) DeleteCampaignId(c ctx.Context, tgId uint64) e.Error {
	err := t.redis.Del(c, campaignQuery(tgId)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}
