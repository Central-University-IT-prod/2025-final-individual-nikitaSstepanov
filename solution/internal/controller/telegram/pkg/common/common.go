package common

import (
	"slices"

	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

type Common struct {
	telegram TelegramUseCase
}

func New(telegram TelegramUseCase) *Common {
	return &Common{
		telegram: telegram,
	}
}

func (c *Common) Cancel(ctx telebot.Context) error {
	ct := ct.GetCtxTg(ctx)

	state, err := c.telegram.GetState(ct, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	if slices.Contains(newStates, state) {
		err = c.telegram.DeleteNew(ct, uint64(ctx.Sender().ID))
		if err != nil {
			return ctx.Send("Что-то пошло не так...")
		}
	}

	err = c.telegram.SetState(ct, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Операция отменена.")
}

func (c *Common) Help(ctx telebot.Context) error {
	return ctx.Send("С помощью AdCore Вы можете управлять рекламными кампаниями, а также смотреть статистику. Вот основные команды бота: \n\n\n1) /start - Начало работы с ботом, вход в систему ▶️\n\n\n2) /new - Создание новой рекламной кампании. Просто вызовите команду и следуйте указаниям 🆕\n\n\n3) /campaigns - Взаимодействие с существующими кампаниями.\n\n\n4) /stats - Просмотр агрегированной статистики по всем кампаниям 📊\n\n\n5) /cancel - Отмена текущей операции. Если вы начали какую-либо операцию, требующую ответа (создание новой кампании, например), вы можете её отменить, вызвав данный команду 🚫 \n\n🛑 Важно: не работает на операции, требующие нажатия кнопки\n\n\n6) /logout - Завершение сессии ↩️")
}
