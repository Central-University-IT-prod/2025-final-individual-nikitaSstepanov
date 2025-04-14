package campaign

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"

	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

type Campaign struct {
	telegram TelegramUseCase
	campaign CampaignUseCase
	ai       AiUseCase
	time     TimeUseCase
}

func New(telegram TelegramUseCase, campaign CampaignUseCase, time TimeUseCase, ai AiUseCase) *Campaign {
	return &Campaign{
		telegram: telegram,
		campaign: campaign,
		ai:       ai,
		time:     time,
	}
}

func (c *Campaign) Campaigns(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil && err.GetCode() != e.NotFound {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	if err != nil {
		return ctx.Send("Вы не вошли в систему. Чтобы начать работу с ботом, вызовите команду /start")
	}

	cnt, err := c.telegram.GetCampaignsCount(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	if cnt == 0 {
		return ctx.Send("У Вас ещё нет рекламных кампаний. Чтобы создать новую кампанию, вызовите функцию /new")
	}

	pagesCount := int(math.Ceil(float64(cnt) / 5))

	campigns, err := c.campaign.Get(ct, id, 5, 1)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	for _, campaign := range campigns {
		row := keyboard.Row(keyboard.Data(campaign.Title, "campaign", campaign.Id, strconv.FormatInt(1, 10), strconv.FormatInt(int64(pagesCount), 10)))

		rows = append(rows, row)
	}

	if pagesCount > 1 {
		row := keyboard.Row(keyboard.Data("➡️", "next", strconv.FormatInt(1, 10), strconv.FormatInt(int64(pagesCount), 10)))

		rows = append(rows, row)
	}

	keyboard.Inline(rows...)

	return ctx.Send("Ваши кампаниии ⬇️", keyboard)
}

func (c *Campaign) Next(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	data := strings.Split(ctx.Callback().Data[1:], "|")

	cur, parseErr := strconv.ParseInt(data[1], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[2], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cur += 1

	campigns, err := c.campaign.Get(ct, id, 5, int(cur))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	for _, campaign := range campigns {
		row := keyboard.Row(keyboard.Data(campaign.Title, "campaign", campaign.Id, strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)))

		rows = append(rows, row)
	}

	if cur < cnt {
		row := keyboard.Row(
			keyboard.Data("⬅️", "prev", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)),
			keyboard.Data("➡️", "next", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)),
		)

		rows = append(rows, row)
	} else {
		prev := keyboard.Row(keyboard.Data("⬅️", "prev", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)))

		rows = append(rows, prev)
	}

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), "Ваши кампаниии ⬇️", keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (c *Campaign) Prev(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	data := strings.Split(ctx.Callback().Data[1:], "|")

	cur, parseErr := strconv.ParseInt(data[1], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[2], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cur -= 1

	campigns, err := c.campaign.Get(ct, id, 5, int(cur))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	for _, campaign := range campigns {
		row := keyboard.Row(keyboard.Data(campaign.Title, "campaign", campaign.Id, strconv.FormatInt(cur, 10), strconv.FormatInt(cnt, 10)))

		rows = append(rows, row)
	}

	if cur == cnt {
		row := keyboard.Row(
			keyboard.Data("⬅️", "prev", strconv.FormatInt(cur, 10), strconv.FormatInt(cnt, 10)),
		)

		rows = append(rows, row)
	} else {
		if cur > 1 {
			row := keyboard.Row(
				keyboard.Data("⬅️", "prev", strconv.FormatInt(cur, 10), strconv.FormatInt(cnt, 10)),
				keyboard.Data("➡️", "next", strconv.FormatInt(cur, 10), strconv.FormatInt(cnt, 10)),
			)

			rows = append(rows, row)
		} else {
			row := keyboard.Row(keyboard.Data("➡️", "next", strconv.FormatInt(cur, 10), strconv.FormatInt(cnt, 10)))

			rows = append(rows, row)
		}
	}

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), "Ваши кампаниии ⬇️", keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (c *Campaign) Campaign(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	cur, parseErr := strconv.ParseInt(data[2], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[3], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("Статистика", "campaign_stats", camId, strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10))),
		keyboard.Row(keyboard.Data("Редактировать", "campaign_update", camId, strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10))),
		keyboard.Row(keyboard.Data("Удалить", "campaign_delete", camId, strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10))),
		keyboard.Row(keyboard.Data("⬅️", "to_list", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10))),
	)

	text := ""

	if len(campaign.Text) > 800 {
		toAdd := campaign.Text[:800] + "......\n\n......Текст слишком длинный для бота. Полный текст вы можете получить через API"
		text = fmt.Sprintf(
			"%s \n\n%s \n\nНачало: %d \nКонец: %d \n\nЛимит показов: %d \nЛимит переходов: %d\n\nСтоимость показа: %g\nСтоимость перехода: %g\n\nТаргетинг: \n\n",
			campaign.Title, toAdd,
			campaign.StartDate, campaign.EndDate,
			campaign.Billing.ImpressionsLimit, campaign.Billing.ClicksLimit,
			campaign.Billing.CostPerImpression, campaign.Billing.CostPerClick,
		)
	} else {
		text = fmt.Sprintf(
			"%s \n\n%s \n\nНачало: %d \nКонец: %d \n\nЛимит показов: %d \nЛимит переходов: %d\n\nСтоимость показа: %g\nСтоимость перехода: %g\n\nТаргетинг: \n\n",
			campaign.Title, campaign.Text,
			campaign.StartDate, campaign.EndDate,
			campaign.Billing.ImpressionsLimit, campaign.Billing.ClicksLimit,
			campaign.Billing.CostPerImpression, campaign.Billing.CostPerClick,
		)
	}

	targeting := campaign.Targeting

	if campaign.Targeting == nil {
		targeting = &entity.Targeting{}
	}

	if targeting.Gender != nil {
		text = fmt.Sprintf("%sПол: %s\n", text, *targeting.Gender)
	} else {
		text = fmt.Sprintf("%sПол: %s\n", text, "не указан")
	}

	if targeting.AgeFrom != nil {
		text = fmt.Sprintf("%sМин. возраст: %d\n", text, *targeting.AgeFrom)
	} else {
		text = fmt.Sprintf("%sМин. возраст: %s\n", text, "не указан")
	}

	if targeting.AgeTo != nil {
		text = fmt.Sprintf("%sМакс. возраст: %d\n", text, *targeting.AgeTo)
	} else {
		text = fmt.Sprintf("%sМакс. возраст: %s\n", text, "не указан")
	}

	if targeting.Location != nil {
		text = fmt.Sprintf("%sЛокация: %s\n", text, *targeting.Location)
	} else {
		text = fmt.Sprintf("%sЛокация: %s\n", text, "не указана")
	}

	buffer, err := c.campaign.DownloadImage(ct, campaign)
	if err != nil {
		_, editErr := bot.Edit(ctx.Message(), text, keyboard)
		if editErr != nil {
			return editErr
		}

		return ctx.Respond()
	}

	reader := bytes.NewReader(buffer)

	photo := &telebot.Photo{File: telebot.FromReader(reader), Caption: text}

	_, editErr := bot.Edit(ctx.Message(), photo, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (c *Campaign) ToList(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	data := strings.Split(ctx.Callback().Data[1:], "|")

	cur, parseErr := strconv.ParseInt(data[1], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[2], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campigns, err := c.campaign.Get(ct, id, 5, int(cur))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	for _, campaign := range campigns {
		row := keyboard.Row(keyboard.Data(campaign.Title, "campaign", campaign.Id, strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)))

		rows = append(rows, row)
	}

	if cnt > 1 {
		if cur != 1 {
			if cur < cnt {
				row := keyboard.Row(
					keyboard.Data("⬅️", "prev", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)),
					keyboard.Data("➡️", "next", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)),
				)

				rows = append(rows, row)
			} else {
				prev := keyboard.Row(keyboard.Data("⬅️", "prev", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)))

				rows = append(rows, prev)
			}
		} else {
			prev := keyboard.Row(keyboard.Data("➡️", "next", strconv.FormatInt(cur, 10), strconv.FormatInt(int64(cnt), 10)))

			rows = append(rows, prev)
		}
	}

	keyboard.Inline(rows...)

	delErr := bot.Delete(ctx.Message())
	if delErr != nil {
		return delErr
	}

	sendErr := ctx.Send("Ваши кампаниии ⬇️", keyboard)
	if sendErr != nil {
		return sendErr
	}

	return ctx.Respond()
}

func (c *Campaign) DailyBillingNext(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	cur, parseErr := strconv.ParseInt(data[4], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[5], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cur += 1

	billing, err := c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := "Статистика по дням:"

	for _, bill := range billing {
		var conversion float32

		if bill.ImpressionsCount == 0 {
			conversion = 0
		} else if bill.ClicksCount == 0 {
			conversion = 0
		} else {
			conversion = float32(bill.ClicksCount) / float32(bill.ImpressionsCount) * 100
		}

		text += fmt.Sprintf(
			"\n\n\nДень %d\n\nКоличество показов: %d\n\nКоличество переходов: %d\n\nКонверсия: %g\n\nПотрачено на показы: %g\n\nПотрачено на переходы: %g\n\nПотрачено всего: %g",
			bill.Date, bill.ImpressionsCount, bill.ClicksCount, conversion, bill.SpentImpressions, bill.SpentClicks,
			bill.SpentImpressions+bill.SpentClicks,
		)
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	if cur < cnt {
		billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur)+1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		nextName := fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date)
		if len(billing) == 1 {
			nextName = fmt.Sprintf("%d ➡️", billing[0].Date)
		}
		billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur)-1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
		rows = append(rows, keyboard.Row(keyboard.Data(prevName, "daily_billing_prev", camId, data[2], data[3], strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10)), keyboard.Data(nextName, "daily_billing_next", camId, data[2], data[3], strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
	} else {
		billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur)-1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
		rows = append(rows, keyboard.Row(keyboard.Data(prevName, "daily_billing_prev", camId, data[2], data[3], strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "campaign_stats", camId, data[2], data[3])))

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), text, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (c *Campaign) DailyBillingPrev(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	cur, parseErr := strconv.ParseInt(data[4], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[5], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cur -= 1

	billing, err := c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := "Статистика по дням:"

	for _, bill := range billing {
		var conversion float32

		if bill.ImpressionsCount == 0 {
			conversion = 0
		} else if bill.ClicksCount == 0 {
			conversion = 0
		} else {
			conversion = float32(bill.ClicksCount) / float32(bill.ImpressionsCount) * 100
		}

		text += fmt.Sprintf(
			"\n\n\nДень %d\n\nКоличество показов: %d\n\nКоличество переходов: %d\n\nКонверсия: %g\n\nПотрачено на показы: %g\n\nПотрачено на переходы: %g\n\nПотрачено всего: %g",
			bill.Date, bill.ImpressionsCount, bill.ClicksCount, conversion, bill.SpentImpressions, bill.SpentClicks,
			bill.SpentImpressions+bill.SpentClicks,
		)
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	if cur < cnt {
		if cur > 1 {
			billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur)+1)
			if err != nil {
				return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
			}
			nextName := fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date)
			if len(billing) == 1 {
				nextName = fmt.Sprintf("%d ➡️", billing[0].Date)
			}
			billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur)-1)
			if err != nil {
				return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
			}
			prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
			rows = append(rows, keyboard.Row(keyboard.Data(prevName, "daily_billing_prev", camId, data[2], data[3], strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10)), keyboard.Data(nextName, "daily_billing_next", camId, data[2], data[3], strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
		} else {
			billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, 2)
			if err != nil {
				return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
			}
			rows = append(rows, keyboard.Row(keyboard.Data(fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date), "daily_billing_next", camId, data[2], data[3], "1", strconv.FormatInt(int64(cnt), 10))))
		}
	} else {
		billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, int(cur)-1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
		rows = append(rows, keyboard.Row(keyboard.Data(prevName, "daily_billing_prev", camId, data[2], data[3], strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "campaign_stats", camId, data[2], data[3])))

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), text, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (c *Campaign) DailyBilling(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	billing, err := c.campaign.CampaignDailyBilling(ct, camId)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt := len(billing)

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	if cnt == 0 {
		rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "campaign_stats", camId, data[2], data[3])))

		keyboard.Inline(rows...)

		_, editErr := bot.Edit(ctx.Message(), "Ежедневная статистика не найдена", keyboard)
		if editErr != nil {
			return editErr
		}

		return ctx.Respond()
	}

	billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, 1)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := "Статистика по дням:"

	for _, bill := range billing {
		var conversion float32

		if bill.ImpressionsCount == 0 {
			conversion = 0
		} else if bill.ClicksCount == 0 {
			conversion = 0
		} else {
			conversion = float32(bill.ClicksCount) / float32(bill.ImpressionsCount) * 100
		}

		text += fmt.Sprintf(
			"\n\n\nДень %d\n\nКоличество показов: %d\n\nКоличество переходов: %d\n\nКонверсия: %g\n\nПотрачено на показы: %g\n\nПотрачено на переходы: %g\n\nПотрачено всего: %g",
			bill.Date, bill.ImpressionsCount, bill.ClicksCount, conversion, bill.SpentImpressions, bill.SpentClicks,
			bill.SpentImpressions+bill.SpentClicks,
		)
	}

	page := int(math.Ceil(float64(cnt) / 3))

	if page > 1 {
		billing, err = c.campaign.CampaignDailyBillingWithPagination(ct, camId, 3, 2)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		rows = append(rows, keyboard.Row(keyboard.Data(fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date), "daily_billing_next", camId, data[2], data[3], "1", strconv.FormatInt(int64(page), 10))))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "campaign_stats", camId, data[2], data[3])))

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), text, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (c *Campaign) Stats(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	billing, err := c.campaign.CampaignBilling(ct, camId)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("По дням", "daily_billing", camId, data[2], data[3])),
		keyboard.Row(keyboard.Data("⬅️", "campaign", camId, data[2], data[3])),
	)

	delErr := bot.Delete(ctx.Message())
	if delErr != nil {
		return delErr
	}

	var conversion float32

	if billing.ImpressionsCount == 0 {
		conversion = 0
	} else if billing.ClicksCount == 0 {
		conversion = 0
	} else {
		conversion = float32(billing.ClicksCount) / float32(billing.ImpressionsCount) * 100
	}

	text := fmt.Sprintf(
		"Количество показов: %d\n\nКоличество переходов: %d\n\nКонверсия: %g\n\nПотрачено на показы: %g\n\nПотрачено на переходы: %g\n\nПотрачено всего: %g",
		billing.ImpressionsCount, billing.ClicksCount, conversion, billing.SpentImpressions, billing.SpentClicks,
		billing.SpentImpressions+billing.SpentClicks,
	)

	sendErr := ctx.Send(text, keyboard)
	if sendErr != nil {
		return sendErr
	}

	return ctx.Respond()
}

func (c *Campaign) Delete(ctx telebot.Context) error {
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("Да", "confirm_delete", camId)),
		keyboard.Row(keyboard.Data("Нет", "campaign", camId, data[2], data[3])),
	)

	_, editErr := bot.EditCaption(ctx.Message(), ctx.Message().Caption+"\n\nВы уверены, что хотите удалить эту кампанию?")
	if editErr != nil {
		_, editErr := bot.Edit(ctx.Message(), ctx.Message().Text+"\n\nВы уверены, что хотите удалить эту кампанию?", keyboard)
		if editErr != nil {
			return editErr
		}
	} else {
		_, editErr = bot.EditReplyMarkup(ctx.Message(), keyboard)
		if editErr != nil {
			return editErr
		}
	}

	return ctx.Respond()
}

func (c *Campaign) ConfirmDelete(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	err = c.campaign.Delete(ct, campaign)
	if err != nil {
		fmt.Println(err)
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	delErr := bot.Delete(ctx.Message())
	if delErr != nil {
		return delErr
	}

	return ctx.Respond()
}

func (c *Campaign) Update(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	camId := data[1]

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	curDay, err := c.time.Get(ct)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{
		keyboard.Row(keyboard.Data("Название", "edit_title", camId), keyboard.Data("Текст", "edit_text", camId)),
		keyboard.Row(keyboard.Data("Изображение", "edit_image", camId)),
		keyboard.Row(keyboard.Data("Стоимость показа", "cost_impression", camId)),
		keyboard.Row(keyboard.Data("Стоимость перехода", "cost_click", camId)),
	}

	if curDay < campaign.StartDate {
		rows = append(rows, keyboard.Row(keyboard.Data("Лимит показов", "limit_impression", camId)))
		rows = append(rows, keyboard.Row(keyboard.Data("Лимит переходов", "limit_click", camId)))
		rows = append(rows, keyboard.Row(keyboard.Data("Начало кампании", "start_date", camId)))
		rows = append(rows, keyboard.Row(keyboard.Data("Конец кампании", "end_date", camId)))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("Пол", "edit_gender", camId), keyboard.Data("Локация", "edit_location", camId)))
	rows = append(rows, keyboard.Row(keyboard.Data("Мин. возраст", "edit_age_from", camId), keyboard.Data("Макс. возраст", "edit_age_to", camId)))

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "campaign", camId, data[2], data[3])))

	keyboard.Inline(rows...)

	_, editErr := bot.EditCaption(ctx.Message(), ctx.Message().Caption+"\n\nВыберите поле для редактирования:")
	if editErr != nil {
		_, editErr := bot.Edit(ctx.Message(), ctx.Message().Text+"\n\nВыберите поле для редактирования:", keyboard)
		if editErr != nil {
			return editErr
		}
	} else {
		_, editErr = bot.EditReplyMarkup(ctx.Message(), keyboard)
		if editErr != nil {
			return editErr
		}
	}

	return ctx.Respond()
}

func (c *Campaign) UpdateAgeTo(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_age_to")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send(`Введите макс. возраст целевой аудитории, либо "-", чтобы оставить поле пустым`)
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingAgeTo(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := ctx.Text()
	targeting := campaign.Targeting

	if text == "-" {
		targeting.AgeTo = nil
	} else {
		age, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
		if parseErr != nil || age < 0 {
			return ctx.Send("Возраст должен быть целым неотрицательным числом")
		}
		toChange := int(age)

		targeting.AgeTo = &toChange
	}

	campaign.Targeting = targeting

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateAgeFrom(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_age_from")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send(`Введите мин. возраст целевой аудитории, либо "-", чтобы оставить поле пустым`)
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingAgeFrom(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := ctx.Text()
	targeting := campaign.Targeting

	if text == "-" {
		targeting.AgeFrom = nil
	} else {
		age, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
		if parseErr != nil || age < 0 {
			return ctx.Send("Возраст должен быть целым неотрицательным числом")
		}
		toChange := int(age)

		targeting.AgeFrom = &toChange
	}

	campaign.Targeting = targeting

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateGender(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_gender")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send(`Введите пол целевой аудитории, либо "-", чтобы оставить поле пустым`)
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingGender(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := ctx.Text()
	targeting := campaign.Targeting
	if text == "-" {
		targeting.Gender = nil
	} else {
		if !slices.Contains([]string{"MALE", "FEMALE", "ALL"}, text) {
			return ctx.Send(`Пол должен принимать одно из трёх значений: "MALE", "FEMALE", "ALL"`)
		}
		gender := types.Gender(text)

		targeting.Gender = &gender
	}

	campaign.Targeting = targeting

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateLocation(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_location")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send(`Введите новую локацию целевой аудитории, либо "-", чтобы оставить поле пустым`)
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingLocation(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	text := ctx.Text()
	targeting := campaign.Targeting
	if text == "-" {
		targeting.Location = nil
	} else {
		targeting.Location = &text
	}

	campaign.Targeting = targeting

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateTitle(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_title")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новое название")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingTitle(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign.Title = ctx.Text()

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) WaitingPrompt(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	sendErr := ctx.Send("Происходит генерация текста ✨ Подождите немного")
	if sendErr != nil {
		return sendErr
	}

	err = c.campaign.Update(ct, campaign, true, ctx.Text())
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (n *Campaign) GenText(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	err := n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_edit_prompt")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	_, editErr := bot.Edit(ctx.Message(), "Введите промпт для генерации текста")
	if editErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) UpdateText(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_text")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("Сгенерировать", "edit_gen_text")),
	)

	sendErr := ctx.Send(`Введите новый текст объявления. Если Вы хотите сгенерировать текст рекламного объявления с помощью нейросети, нажмите на кнопку "Сгенерировать"`, keyboard)
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingText(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign.Text = ctx.Text()

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateImage(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_image")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send(`Отправьте новое изображение или "-", чтобы удалить изображение`)
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingImageText(ctx telebot.Context) error {
	text := ctx.Text()

	if text == "-" {
		ct := ct.GetCtxTg(ctx)

		id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}

		camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}

		campaign := &entity.Campaign{
			Id:           camId,
			AdvertiserId: id,
		}

		err = c.campaign.DeleteImage(ct, campaign)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}

		err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}

		err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}

		return ctx.Send("Кампания успешно обновлена")
	} else {
		return ctx.Send(`Отправьте новое изображение или "-", чтобы удалить изображение`)
	}
}

func (c *Campaign) WaitingImage(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	photo := ctx.Message().Photo

	file, err := bot.FileByID(photo.FileID)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	fileUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)

	resp, err := http.Get(fileUrl)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}
	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	image := &entity.Image{
		Name:        file.FilePath,
		Buffer:      buffer,
		Size:        photo.FileSize,
		ContentType: "image/jpg",
	}

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	_, err = c.campaign.UploadImage(ct, campaign, image)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateCostImpression(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_cost_impression")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новую стоимость показа")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingCostImpression(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cost, parseErr := strconv.ParseFloat(ctx.Text(), 64)
	if parseErr != nil || cost < 0 {
		return ctx.Send("Стоимось показа должна быть неотрицательным числом")
	}

	campaign.Billing.CostPerImpression = float32(cost)

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateCostClick(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_cost_click")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новую стоимость перехода")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingCostClick(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cost, parseErr := strconv.ParseFloat(ctx.Text(), 64)
	if parseErr != nil || cost < 0 {
		return ctx.Send("Стоимось перехода должна быть неотрицательным числом")
	}

	campaign.Billing.CostPerClick = float32(cost)

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateImpressionsLimit(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_limit_impressions")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новый лимит показов")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingImpressionLimit(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	limit, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || limit < 0 {
		return ctx.Send("Лимит показов должен быть целым неоотрицательным числом")
	}

	campaign.Billing.ImpressionsLimit = int(limit)

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateClickLimit(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_limit_click")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новый лимит переходов")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingClickLimit(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	limit, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || limit < 0 {
		return ctx.Send("Лимит переходов должен быть целым неоотрицательным числом")
	}

	campaign.Billing.ClicksLimit = int(limit)

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateStartDate(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_start_date")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новую дату начала кампании")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingStartDate(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	date, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || date < 0 {
		return ctx.Send("Дата начала кампании должна быть челым неотрицательным числом")
	}

	campaign.StartDate = int(date)

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}

func (c *Campaign) UpdateEndDate(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")
	id := data[1]

	err := c.telegram.SetCampaignId(ct, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "waiting_edit_end_date")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	bot.Delete(ctx.Message())

	sendErr := ctx.Send("Введите новую дату конца кампании")
	if sendErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Respond()
}

func (c *Campaign) WaitingEndDate(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	id, err := c.telegram.GetSession(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	camId, err := c.telegram.GetCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err = c.campaign.GetById(ct, campaign)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	date, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || date < 0 {
		return ctx.Send("Дата конца кампании должна быть челым неотрицательным числом")
	}

	campaign.EndDate = int(date)

	err = c.campaign.Update(ct, campaign, false, "")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.DeleteCampaignId(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	return ctx.Send("Кампания успешно обновлена")
}
