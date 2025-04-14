package new

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

type NewCampaign struct {
	telegram TelegramUseCase
	campaign CampaignUseCase
	ai       AiUseCase
}

func New(telegram TelegramUseCase, campain CampaignUseCase, ai AiUseCase) *NewCampaign {
	return &NewCampaign{
		telegram: telegram,
		campaign: campain,
		ai:       ai,
	}
}

func (n *NewCampaign) New(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	session, err := n.telegram.GetSession(c, uint64(ctx.Sender().ID))
	if err != nil && err.GetCode() != e.NotFound {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	if err != nil {
		return ctx.Send("Вы не вошли в систему. Чтобы начать работу с ботом, вызовите команду /start")
	}

	data := &entity.CampaignData{
		AdvertiserId: session,
	}

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_title")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите название рекламной кампании")
}

func (n *NewCampaign) WaitingImage(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
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

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	createErr := n.createNew(c, data, image)
	if createErr != nil && createErr.GetCode() != e.BadInput {
		fmt.Println(createErr)
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	if createErr != nil {
		return ctx.Send("Вы ввели некорректные данные. Попробуйте создать кампанию заново.")
	}

	err = n.telegram.DeleteNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так...")
	}

	setErr := n.telegram.SetState(c, uint64(ctx.Sender().ID), "nothing")
	if setErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Кампания успешно создана 🎉")
}

func (n *NewCampaign) WaitingImageText(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	if ctx.Text() == "-" {
		data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		err = n.createNew(c, data, nil)
		if err != nil && err.GetCode() != e.BadInput {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		if err != nil {
			return ctx.Send("Вы ввели некорректные данные. Попробуйте создать кампанию заново.")
		}

		err = n.telegram.DeleteNew(c, uint64(ctx.Sender().ID))
		if err != nil {
			return ctx.Send("Что-то пошло не так...")
		}

		err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "nothing")
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		return ctx.Send("Кампания успешно создана 🎉")
	} else {
		return ctx.Send(`Отправьте изображение, которое будет отображаться при показе объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
	}
}

func (n *NewCampaign) WaitingLocation(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	text := ctx.Text()

	if text == "-" {
		err := n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_image")
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		return ctx.Send(`Отправьте изображение, которое будет отображаться при показе объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
	}

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	data.Location = &text

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_image")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send(`Отправьте изображение, которое будет отображаться при показе объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
}

func (n *NewCampaign) WaitingAgeTo(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	text := ctx.Text()

	if text == "-" {
		err := n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_location")
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		return ctx.Send(`Введите локацию целевой аудитории объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
	}

	age, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || age < 0 {
		return ctx.Send("Возраст должен быть целым неотрицательным числом")
	}

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	to := int(age)

	data.AgeFrom = &to

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_location")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send(`Введите локацию целевой аудитории объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
}

func (n *NewCampaign) WaitingAgeFrom(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	text := ctx.Text()

	if text == "-" {
		err := n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_age_to")
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		return ctx.Send(`Введите максимальный возраст среди целевой аудитории объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
	}

	age, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || age < 0 {
		return ctx.Send("Возраст должен быть целым неотрицательным числом")
	}

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	from := int(age)

	data.AgeFrom = &from

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_age_to")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send(`Введите максимальный возраст среди целевой аудитории объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
}

func (n *NewCampaign) WaitingGender(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	text := ctx.Text()

	if text == "-" {
		err := n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_age_from")
		if err != nil {
			return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
		}

		return ctx.Send(`Введите минимальный возраст среди целевой аудитории объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
	}

	if !slices.Contains([]string{"MALE", "FEMALE", "ALL"}, text) {
		return ctx.Send(`Пол должен принимать одно из трёх значений: "MALE", "FEMALE", "ALL"`)
	}

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	gender := types.Gender(text)

	data.Gender = &gender

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_age_from")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send(`Введите минимальный возраст среди целевой аудитории объявления. Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
}

func (n *NewCampaign) WaitingClickCost(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	cost, parseErr := strconv.ParseFloat(ctx.Text(), 64)
	if parseErr != nil || cost < 0 {
		return ctx.Send("Стоимось перехода должна быть неотрицательным числом")
	}

	data.CostPerClick = float32(cost)

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_gender")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send(`Введите пол целевой аудитории объявления (одно из трёх значений: "MALE", "FEMALE", "ALL"). Если Вы хотите оставить данное поле пустым, введите прочерк ("-")`)
}

func (n *NewCampaign) WaitingImpressionCost(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	cost, parseErr := strconv.ParseFloat(ctx.Text(), 64)
	if parseErr != nil || cost < 0 {
		return ctx.Send("Стоимось показа должна быть неотрицательным числом")
	}

	data.CostPerImpression = float32(cost)

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_click_cost")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите стоимость перехода по объявлению")
}

func (n *NewCampaign) WaitingClicksLimit(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	limit, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || limit < 0 {
		return ctx.Send("Лимит переходов должен быть целым неоотрицательным числом")
	}

	data.ClicksLimit = int(limit)

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_impression_cost")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите стоимость показа объявления")
}

func (n *NewCampaign) WaitingImpressionsLimit(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	limit, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || limit < 0 {
		return ctx.Send("Лимит показов должен быть целым неоотрицательным числом")
	}

	data.ImpressionsLimit = int(limit)

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_clicks_limit")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите лимит переходов по объявлению")
}

func (n *NewCampaign) WaitingEndDate(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	date, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || date < 0 {
		return ctx.Send("Дата конца кампании должна быть челым неотрицательным числом")
	}

	data.EndDate = int(date)

	if data.EndDate < data.StartDate {
		return ctx.Send("Дата конца не может быть раньше даты начала")
	}

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_impressions_limit")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите лимит показов объявления")
}

func (n *NewCampaign) WaitingStartDate(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	date, parseErr := strconv.ParseInt(ctx.Text(), 10, 64)
	if parseErr != nil || date < 0 {
		return ctx.Send("Дата начала кампании должна быть челым неотрицательным числом")
	}

	data.StartDate = int(date)

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_end_date")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите дату конца кампании")
}

func (n *NewCampaign) WaitingText(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	data.Text = ctx.Text()

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_start_date")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите дату начала кампании")
}

func (n *NewCampaign) WaitingPrompt(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	sendErr := ctx.Send("Происходит генерация текста ✨ Подождите немного")
	if sendErr != nil {
		return sendErr
	}

	text, err := n.ai.GenText(ctx.Text())
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	data.Text = text

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_start_date")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Send("Введите дату начала кампании")
}

func (n *NewCampaign) GenText(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)
	bot := ctx.Get("bot").(*telebot.Bot)

	err := n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_prompt")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	_, editErr := bot.Edit(ctx.Message(), "Введите промпт для генерации текста")
	if editErr != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	return ctx.Respond()
}

func (n *NewCampaign) WaitingTitle(ctx telebot.Context) error {
	c := ct.GetCtxTg(ctx)

	data, err := n.telegram.GetNew(c, uint64(ctx.Sender().ID))
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	data.Title = ctx.Text()

	err = n.telegram.SetNew(c, uint64(ctx.Sender().ID), data)
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	err = n.telegram.SetState(c, uint64(ctx.Sender().ID), "waiting_text")
	if err != nil {
		return ctx.Send("Что-то пошло не так... Попробуйте ешё раз.")
	}

	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		keyboard.Row(keyboard.Data("Сгенерировать", "new_gen_text")),
	)

	return ctx.Send(`Введите текст объявления. Если Вы хотите сгенерировать текст рекламного объявления с помощью нейросети, нажмите на кнопку "Сгенерировать"`, keyboard)
}

func (n *NewCampaign) createNew(c ctx.Context, data *entity.CampaignData, image *entity.Image) e.Error {
	campaign := &entity.Campaign{
		AdvertiserId: data.AdvertiserId,
		Title:        data.Title,
		Text:         data.Text,
		StartDate:    data.StartDate,
		EndDate:      data.EndDate,
		Billing: &entity.Billing{
			ImpressionsLimit:  data.ImpressionsLimit,
			ClicksLimit:       data.ClicksLimit,
			CostPerImpression: data.CostPerImpression,
			CostPerClick:      data.CostPerClick,
		},
		Targeting: &entity.Targeting{
			Gender:   data.Gender,
			AgeFrom:  data.AgeFrom,
			AgeTo:    data.AgeTo,
			Location: data.Location,
		},
	}

	err := n.campaign.Create(c, campaign, false, "")
	if err != nil {
		return err
	}

	if image != nil {
		_, err := n.campaign.UploadImage(c, campaign, image)
		if err != nil {
			return err
		}
	}

	return nil
}
