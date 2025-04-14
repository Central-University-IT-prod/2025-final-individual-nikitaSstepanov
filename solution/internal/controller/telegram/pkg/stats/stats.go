package stats

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

type Stats struct {
	telegram TelegramUseCase
	campaign CampaignUseCase
}

func New(telegram TelegramUseCase, campaign CampaignUseCase) *Stats {
	return &Stats{
		telegram: telegram,
		campaign: campaign,
	}
}

func (s *Stats) Get(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	id, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	billing, err := s.campaign.AdvertiserBilling(c, id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	var conversion float32

	if billing.ImpressionsCount == 0 {
		conversion = 0
	} else if billing.ClicksCount == 0 {
		conversion = 0
	} else {
		conversion = float32(billing.ClicksCount) / float32(billing.ImpressionsCount) * 100
	}

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("По дням", "stats_daily")),
	)

	text := fmt.Sprintf(
		"Количество показов: %d\n\nКоличество переходов: %d\n\nКонверсия: %g\n\nПотрачено на показы: %g\n\nПотрачено на переходы: %g\n\nПотрачено всего: %g",
		billing.ImpressionsCount, billing.ClicksCount, conversion, billing.SpentImpressions, billing.SpentClicks,
		billing.SpentImpressions+billing.SpentClicks,
	)

	return ctx.Send(text, keyboard)
}

func (s *Stats) GetCallback(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	billing, err := s.campaign.AdvertiserBilling(c, id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	var conversion float32

	if billing.ImpressionsCount == 0 {
		conversion = 0
	} else if billing.ClicksCount == 0 {
		conversion = 0
	} else {
		conversion = float32(billing.ClicksCount) / float32(billing.ImpressionsCount) * 100
	}

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("По дням", "stats_daily")),
	)

	text := fmt.Sprintf(
		"Количество показов: %d\n\nКоличество переходов: %d\n\nКонверсия: %g\n\nПотрачено на показы: %g\n\nПотрачено на переходы: %g\n\nПотрачено всего: %g",
		billing.ImpressionsCount, billing.ClicksCount, conversion, billing.SpentImpressions, billing.SpentClicks,
		billing.SpentImpressions+billing.SpentClicks,
	)

	bot.Edit(ctx.Message(), text, keyboard)

	return ctx.Respond()
}

func (s *Stats) GetDaily(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	billing, err := s.campaign.AdvertiserDailyBill(c, id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt := len(billing)
	fmt.Println(cnt)
	keyboard := &telebot.ReplyMarkup{}

	rows := []telebot.Row{}

	if cnt == 0 {
		rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "stats_callback")))

		keyboard.Inline(rows...)

		_, editErr := bot.Edit(ctx.Message(), "Ежедневная статистика не найдена", keyboard)
		if editErr != nil {
			return editErr
		}

		return ctx.Respond()
	}

	billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, 1)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}
	fmt.Println(billing)
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
		billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, 2)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		rows = append(rows, keyboard.Row(keyboard.Data(fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date), "stats_daily_next", "1", strconv.FormatInt(int64(page), 10))))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "stats_callback")))

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), text, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (s *Stats) GetDailyNext(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	id, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
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

	billing, err := s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur))
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
		billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur)+1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		nextName := fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date)
		if len(billing) == 1 {
			nextName = fmt.Sprintf("%d ➡️", billing[0].Date)
		}
		billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur)-1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
		rows = append(rows, keyboard.Row(keyboard.Data(prevName, "stats_daily_prev", strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10)), keyboard.Data(nextName, "stats_daily_next", strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
	} else {
		billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur)-1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
		rows = append(rows, keyboard.Row(keyboard.Data(prevName, "stats_daily_prev", strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "stats_callback")))

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), text, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}

func (s *Stats) GetDailyPrev(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	data := strings.Split(ctx.Callback().Data[1:], "|")

	id, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cur, parseErr := strconv.ParseInt(data[1], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cnt, parseErr := strconv.ParseInt(data[2], 10, 64)
	if parseErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
	}

	cur -= 1

	billing, err := s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur))
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
			billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur)+1)
			if err != nil {
				return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
			}
			nextName := fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date)
			if len(billing) == 1 {
				nextName = fmt.Sprintf("%d ➡️", billing[0].Date)
			}
			billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur)-1)
			if err != nil {
				return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
			}
			prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
			rows = append(rows, keyboard.Row(keyboard.Data(prevName, "stats_daily_prev", strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10)), keyboard.Data(nextName, "stats_daily_next", strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
		} else {
			billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, 2)
			if err != nil {
				return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
			}
			rows = append(rows, keyboard.Row(keyboard.Data(fmt.Sprintf("%d-%d ➡️", billing[0].Date, billing[len(billing)-1].Date), "stats_daily_next", "1", strconv.FormatInt(int64(cnt), 10))))
		}
	} else {
		billing, err = s.campaign.AdvertiserDailyBillWithPagination(c, id, 3, int(cur)-1)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ещё раз.")
		}
		prevName := fmt.Sprintf("⬅️ %d-%d", billing[0].Date, billing[len(billing)-1].Date)
		rows = append(rows, keyboard.Row(keyboard.Data(prevName, "stats_daily_prev", strconv.FormatInt(int64(cur), 10), strconv.FormatInt(int64(cnt), 10))))
	}

	rows = append(rows, keyboard.Row(keyboard.Data("⬅️", "stats_callback")))

	keyboard.Inline(rows...)

	_, editErr := bot.Edit(ctx.Message(), text, keyboard)
	if editErr != nil {
		return editErr
	}

	return ctx.Respond()
}
