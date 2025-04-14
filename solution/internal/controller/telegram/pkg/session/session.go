package session

import (
	"fmt"

	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram/validator"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

type Session struct {
	telegram   TelegramUseCase
	advertiser AdvertiserUseCase
}

func New(telegram TelegramUseCase, advertiser AdvertiserUseCase) *Session {
	return &Session{
		telegram:   telegram,
		advertiser: advertiser,
	}
}

func (s *Session) Start(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	session, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil && err.GetCode() != e.NotFound {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	if err == nil {
		advertiser, err := s.advertiser.GetById(c, session)
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		return ctx.Send(fmt.Sprintf("Вы уже вошли в систему как %s. Чтобы ознакомиться с функционалом бота, вызовите команду /help", advertiser.Name))
	}

	err = s.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_id")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Здравствуйте! Чтобы продолжить работу, отправьте свой ID рекламодателя.")
}

func (s *Session) WaitingId(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	id := ctx.Text()

	if err := validator.UUID(id); err != nil {
		return ctx.Send("ID должен иметь формат UUID. Попробуйте ещё раз.")
	}

	_, err := s.advertiser.GetById(c, id)
	if err != nil {
		return ctx.Send("Рекламодатель с таким ID не найден 😔")
	}

	err = s.telegram.SetSession(c, uint64(ctx.Sender().ID), id)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = s.telegram.SetState(c, uint64(ctx.Sender().ID), "nothing")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Отлично, вы вошли в систему 😉 Чтобы ознакомиться с функционалом бота, вызовите команду /help")
}

func (s *Session) Logout(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	_, err := s.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil && err.GetCode() != e.NotFound {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	if err != nil {
		return ctx.Send("Вы не вошли в систему. Чтобы начать работу с ботом, вызовите команду /start")
	}

	err = s.telegram.DeleteSession(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Сессия успешно завершена. Чтобы начать новую, вызовите команду /start")
}
