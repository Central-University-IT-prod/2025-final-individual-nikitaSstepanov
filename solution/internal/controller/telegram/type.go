package telegram

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gopkg.in/telebot.v4"
)

type SessionHandler interface {
	Start(ctx telebot.Context) error
	WaitingId(ctx telebot.Context) error
	Logout(ctx telebot.Context) error
}

type NewHandler interface {
	WaitingPrompt(ctx telebot.Context) error
	GenText(ctx telebot.Context) error
	New(ctx telebot.Context) error
	WaitingImage(ctx telebot.Context) error
	WaitingImageText(ctx telebot.Context) error
	WaitingLocation(ctx telebot.Context) error
	WaitingAgeTo(ctx telebot.Context) error
	WaitingAgeFrom(ctx telebot.Context) error
	WaitingGender(ctx telebot.Context) error
	WaitingClickCost(ctx telebot.Context) error
	WaitingImpressionCost(ctx telebot.Context) error
	WaitingClicksLimit(ctx telebot.Context) error
	WaitingImpressionsLimit(ctx telebot.Context) error
	WaitingEndDate(ctx telebot.Context) error
	WaitingStartDate(ctx telebot.Context) error
	WaitingText(ctx telebot.Context) error
	WaitingTitle(ctx telebot.Context) error
}

type StatsHandler interface {
	Get(ctx telebot.Context) error
	GetCallback(ctx telebot.Context) error
	GetDaily(ctx telebot.Context) error
	GetDailyNext(ctx telebot.Context) error
	GetDailyPrev(ctx telebot.Context) error
}

type CampaignHandler interface {
	WaitingPrompt(ctx telebot.Context) error
	GenText(ctx telebot.Context) error
	Update(ctx telebot.Context) error
	ConfirmDelete(ctx telebot.Context) error
	Delete(ctx telebot.Context) error
	ToList(ctx telebot.Context) error
	Campaign(ctx telebot.Context) error
	Campaigns(ctx telebot.Context) error
	Next(ctx telebot.Context) error
	Prev(ctx telebot.Context) error
	UpdateTitle(ctx telebot.Context) error
	UpdateText(ctx telebot.Context) error
	UpdateImage(ctx telebot.Context) error
	UpdateCostImpression(ctx telebot.Context) error
	UpdateCostClick(ctx telebot.Context) error
	UpdateImpressionsLimit(ctx telebot.Context) error
	UpdateClickLimit(ctx telebot.Context) error
	UpdateStartDate(ctx telebot.Context) error
	UpdateEndDate(ctx telebot.Context) error
	WaitingText(ctx telebot.Context) error
	WaitingTitle(ctx telebot.Context) error
	WaitingImageText(ctx telebot.Context) error
	WaitingImage(ctx telebot.Context) error
	WaitingCostImpression(ctx telebot.Context) error
	WaitingCostClick(ctx telebot.Context) error
	WaitingImpressionLimit(ctx telebot.Context) error
	WaitingClickLimit(ctx telebot.Context) error
	WaitingStartDate(ctx telebot.Context) error
	WaitingEndDate(ctx telebot.Context) error
	UpdateLocation(ctx telebot.Context) error
	WaitingLocation(ctx telebot.Context) error
	UpdateGender(ctx telebot.Context) error
	WaitingGender(ctx telebot.Context) error
	UpdateAgeTo(ctx telebot.Context) error
	WaitingAgeFrom(ctx telebot.Context) error
	UpdateAgeFrom(ctx telebot.Context) error
	WaitingAgeTo(ctx telebot.Context) error
	Stats(ctx telebot.Context) error
	DailyBilling(ctx telebot.Context) error
	DailyBillingNext(ctx telebot.Context) error
	DailyBillingPrev(ctx telebot.Context) error
}

type TelegramUseCase interface {
	GetState(c ctx.Context, id uint64) (string, e.Error)
	SetState(c ctx.Context, id uint64, state string) e.Error
	GetSession(c ctx.Context, tgId uint64) (string, e.Error)
	SetSession(c ctx.Context, tgId uint64, id string) e.Error
	DeleteSession(c ctx.Context, tgId uint64) e.Error
}

type CommonHandler interface {
	Cancel(ctx telebot.Context) error
	Help(ctx telebot.Context) error
}

type Middleware interface {
	InitLogger(c ctx.Context) telebot.MiddlewareFunc
}
