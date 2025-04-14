package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/middleware"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/pkg/campaign"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/pkg/common"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/pkg/new"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/pkg/session"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/pkg/stats"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

type Bot struct {
	session  SessionHandler
	new      NewHandler
	common   CommonHandler
	campaign CampaignHandler
	mid      Middleware
	telegram TelegramUseCase
	stats    StatsHandler
	cfg      *Config
}

type Config struct {
	Token   string        `env:"TG_TOKEN"`
	Timeout time.Duration `yaml:"timeout"`
}

func New(uc *usecase.UseCase, cfg *Config) *Bot {
	return &Bot{
		session:  session.New(uc.Telegram, uc.Advertiser),
		new:      new.New(uc.Telegram, uc.Campaign, uc.Ai),
		common:   common.New(uc.Telegram),
		campaign: campaign.New(uc.Telegram, uc.Campaign, uc.Time, uc.Ai),
		stats:    stats.New(uc.Telegram, uc.Campaign),
		telegram: uc.Telegram,
		mid:      middleware.New(),
		cfg:      cfg,
	}
}

func (b *Bot) InitBot(c ctx.Context) *telebot.Bot {
	pref := telebot.Settings{
		Token: b.cfg.Token,
		Poller: &telebot.LongPoller{
			Timeout: b.cfg.Timeout,
		},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log := c.Logger()

		log.Error("Can`t init telegram bot.", sl.ErrAttr(err))
		panic("App start error.")
	}

	bot.Use(b.mid.InitLogger(c))

	bot.Handle("/cancel", b.common.Cancel)

	bot.Handle(telebot.OnCallback, func(ctx telebot.Context) error {
		mode := strings.Split(ctx.Callback().Data[1:], "|")[0]

		switch mode {
		case "next":
			ctx.Set("bot", bot)
			return b.campaign.Next(ctx)
		case "prev":
			ctx.Set("bot", bot)
			return b.campaign.Prev(ctx)
		case "to_list":
			ctx.Set("bot", bot)
			return b.campaign.ToList(ctx)
		case "campaign_delete":
			ctx.Set("bot", bot)
			return b.campaign.Delete(ctx)
		case "confirm_delete":
			ctx.Set("bot", bot)
			return b.campaign.ConfirmDelete(ctx)
		case "campaign_update":
			ctx.Set("bot", bot)
			return b.campaign.Update(ctx)
		case "edit_title":
			ctx.Set("bot", bot)
			return b.campaign.UpdateTitle(ctx)
		case "edit_text":
			ctx.Set("bot", bot)
			return b.campaign.UpdateText(ctx)
		case "edit_image":
			ctx.Set("bot", bot)
			return b.campaign.UpdateImage(ctx)
		case "cost_impression":
			ctx.Set("bot", bot)
			return b.campaign.UpdateCostImpression(ctx)
		case "cost_click":
			ctx.Set("bot", bot)
			return b.campaign.UpdateCostClick(ctx)
		case "limit_impression":
			ctx.Set("bot", bot)
			return b.campaign.UpdateImpressionsLimit(ctx)
		case "limit_click":
			ctx.Set("bot", bot)
			return b.campaign.UpdateClickLimit(ctx)
		case "start_date":
			ctx.Set("bot", bot)
			return b.campaign.UpdateStartDate(ctx)
		case "end_date":
			ctx.Set("bot", bot)
			return b.campaign.UpdateEndDate(ctx)
		case "edit_location":
			ctx.Set("bot", bot)
			return b.campaign.UpdateLocation(ctx)
		case "edit_gender":
			ctx.Set("bot", bot)
			return b.campaign.UpdateGender(ctx)
		case "edit_age_from":
			ctx.Set("bot", bot)
			return b.campaign.UpdateAgeFrom(ctx)
		case "edit_age_to":
			ctx.Set("bot", bot)
			return b.campaign.UpdateAgeTo(ctx)
		case "campaign_stats":
			ctx.Set("bot", bot)
			return b.campaign.Stats(ctx)
		case "daily_billing":
			ctx.Set("bot", bot)
			return b.campaign.DailyBilling(ctx)
		case "daily_billing_next":
			ctx.Set("bot", bot)
			return b.campaign.DailyBillingNext(ctx)
		case "daily_billing_prev":
			ctx.Set("bot", bot)
			return b.campaign.DailyBillingPrev(ctx)
		case "stats_callback":
			ctx.Set("bot", bot)
			return b.stats.GetCallback(ctx)
		case "stats_daily":
			ctx.Set("bot", bot)
			return b.stats.GetDaily(ctx)
		case "stats_daily_next":
			ctx.Set("bot", bot)
			return b.stats.GetDailyNext(ctx)
		case "stats_daily_prev":
			ctx.Set("bot", bot)
			return b.stats.GetDailyPrev(ctx)
		case "new_gen_text":
			ctx.Set("bot", bot)
			return b.new.GenText(ctx)
		case "campaign":
			ctx.Set("bot", bot)
			return b.campaign.Campaign(ctx)
		case "edit_gen_text":
			ctx.Set("bot", bot)
			return b.campaign.GenText(ctx)
		}

		return ctx.Respond()
	})

	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		c := ct.GetCtxTg(ctx)

		state, err := b.telegram.GetState(c, uint64(ctx.Sender().ID))
		if err != nil && err.GetCode() != e.NotFound {
			return ctx.Send("Что-то пошло не так...")
		}

		if err != nil {
			state = "nothing"
		}

		switch state {

		case "nothing":
			switch ctx.Text() {
			case "/start":
				return b.session.Start(ctx)
			case "/logout":
				return b.session.Logout(ctx)
			case "/new":
				return b.new.New(ctx)
			case "/campaigns":
				return b.campaign.Campaigns(ctx)
			case "/stats":
				return b.stats.Get(ctx)
			default:
				return b.common.Help(ctx)
			}

		case "waiting_id":
			return b.session.WaitingId(ctx)

		case "waiting_title":
			return b.new.WaitingTitle(ctx)

		case "waiting_text":
			return b.new.WaitingText(ctx)

		case "waiting_start_date":
			return b.new.WaitingStartDate(ctx)

		case "waiting_end_date":
			return b.new.WaitingEndDate(ctx)

		case "waiting_impressions_limit":
			return b.new.WaitingImpressionsLimit(ctx)

		case "waiting_clicks_limit":
			return b.new.WaitingClicksLimit(ctx)

		case "waiting_impression_cost":
			return b.new.WaitingImpressionCost(ctx)

		case "waiting_click_cost":
			return b.new.WaitingClickCost(ctx)

		case "waiting_gender":
			return b.new.WaitingGender(ctx)

		case "waiting_age_from":
			return b.new.WaitingAgeFrom(ctx)

		case "waiting_age_to":
			return b.new.WaitingAgeTo(ctx)

		case "waiting_location":
			return b.new.WaitingLocation(ctx)

		case "waiting_image":
			return b.new.WaitingImageText(ctx)

		case "waiting_edit_title":
			return b.campaign.WaitingTitle(ctx)

		case "waiting_edit_text":
			return b.campaign.WaitingText(ctx)

		case "waiting_edit_image":
			return b.campaign.WaitingImageText(ctx)

		case "waiting_edit_cost_impression":
			return b.campaign.WaitingCostImpression(ctx)

		case "waiting_edit_cost_click":
			return b.campaign.WaitingCostClick(ctx)

		case "waiting_edit_limit_impressions":
			return b.campaign.WaitingImpressionLimit(ctx)

		case "waiting_edit_limit_click":
			return b.campaign.WaitingClickLimit(ctx)

		case "waiting_edit_start_date":
			return b.campaign.WaitingStartDate(ctx)

		case "waiting_edit_end_date":
			return b.campaign.WaitingEndDate(ctx)

		case "waiting_edit_location":
			return b.campaign.WaitingLocation(ctx)

		case "waiting_edit_gender":
			return b.campaign.WaitingGender(ctx)

		case "waiting_edit_age_to":
			return b.campaign.WaitingAgeTo(ctx)

		case "waiting_edit_age_from":
			return b.campaign.WaitingAgeFrom(ctx)

		case "waiting_prompt":
			return b.new.WaitingPrompt(ctx)

		case "waiting_edit_prompt":
			return b.campaign.WaitingPrompt(ctx)

		}

		return nil
	})

	bot.Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
		c := ct.GetCtxTg(ctx)

		state, err := b.telegram.GetState(c, uint64(ctx.Sender().ID))
		if err != nil && err.GetCode() != e.NotFound {
			fmt.Println(err)
			return ctx.Send("Что-то пошло не так...")
		}

		if err != nil {
			state = "nothing"
		}

		switch state {
		case "nothing":
			return b.common.Help(ctx)

		case "waiting_image":
			ctx.Set("bot", bot)
			return b.new.WaitingImage(ctx)

		case "waiting_edit_image":
			ctx.Set("bot", bot)
			return b.campaign.WaitingImage(ctx)
		}

		return nil
	})

	return bot
}
