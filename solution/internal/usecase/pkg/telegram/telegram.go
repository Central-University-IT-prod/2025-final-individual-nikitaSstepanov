package telegram

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Telegram struct {
	telegram TelegramStorage
	campaign CampaignStorage
}

func New(telegram TelegramStorage, campaign CampaignStorage) *Telegram {
	return &Telegram{
		telegram: telegram,
		campaign: campaign,
	}
}

func (t *Telegram) GetState(c ctx.Context, id uint64) (string, e.Error) {
	return t.telegram.GetState(c, id)
}

func (t *Telegram) SetState(c ctx.Context, id uint64, state string) e.Error {
	return t.telegram.SetState(c, id, state)
}

func (t *Telegram) GetSession(c ctx.Context, tgId uint64) (string, e.Error) {
	return t.telegram.GetSession(c, tgId)
}

func (t *Telegram) SetSession(c ctx.Context, tgId uint64, id string) e.Error {
	return t.telegram.SetSession(c, tgId, id)
}

func (t *Telegram) DeleteSession(c ctx.Context, tgId uint64) e.Error {
	return t.telegram.DeleteSession(c, tgId)
}

func (t *Telegram) GetNew(c ctx.Context, tgId uint64) (*entity.CampaignData, e.Error) {
	return t.telegram.GetNew(c, tgId)
}

func (t *Telegram) SetNew(c ctx.Context, tgId uint64, data *entity.CampaignData) e.Error {
	return t.telegram.SetNew(c, tgId, data)
}

func (t *Telegram) DeleteNew(c ctx.Context, tgId uint64) e.Error {
	return t.telegram.DeleteNew(c, tgId)
}

func (t *Telegram) GetCampaignsCount(c ctx.Context, tgId uint64) (int, e.Error) {
	id, err := t.GetSession(c, tgId)
	if err != nil {
		return 0, err
	}

	campaigns, err := t.campaign.Get(c, id)
	if err != nil && err.GetCode() != e.NotFound {
		return 0, err
	}

	if err != nil {
		return 0, nil
	}

	return len(campaigns), nil
}

func (t *Telegram) GetCampaignId(c ctx.Context, tgId uint64) (string, e.Error) {
	return t.telegram.GetCampaignId(c, tgId)
}

func (t *Telegram) SetCampaignId(c ctx.Context, tgId uint64, id string) e.Error {
	return t.telegram.SetCampaignId(c, tgId, id)
}

func (t *Telegram) DeleteCampaignId(c ctx.Context, tgId uint64) e.Error {
	return t.telegram.DeleteCampaignId(c, tgId)
}
