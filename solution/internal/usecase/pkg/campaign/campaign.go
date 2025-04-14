package campaign

import (
	"strings"

	"github.com/google/uuid"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Campaign struct {
	campaign   CampaignStorage
	advertiser AdvertiserStorage
	client     ClientStorage
	score      ScoreStorage
	image      ImageStorage
	time       TimeUseCase
	moderate   ModerateUseCase
	ai         AiUseCase
}

func New(campaign CampaignStorage, advertiser AdvertiserStorage, time TimeUseCase, client ClientStorage, image ImageStorage, score ScoreStorage, ai AiUseCase, moderate ModerateUseCase) *Campaign {
	return &Campaign{
		campaign:   campaign,
		advertiser: advertiser,
		client:     client,
		time:       time,
		image:      image,
		score:      score,
		ai:         ai,
		moderate:   moderate,
	}
}

func (c *Campaign) GetById(ctx ctx.Context, campaign *entity.Campaign) (*entity.Campaign, e.Error) {
	advert, err := c.advertiser.GetById(ctx, campaign.AdvertiserId)
	if err != nil {
		return nil, err
	}

	campaign, err = c.campaign.GetById(ctx, campaign.Id)
	if err != nil {
		return nil, err
	}

	if campaign.AdvertiserId != advert.Id {
		return nil, forbiddenErr
	}

	return campaign, nil
}

func (c *Campaign) Get(ctx ctx.Context, advertiserId string, size int, page int) ([]*entity.Campaign, e.Error) {
	_, err := c.advertiser.GetById(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	return c.campaign.GetWithPagination(ctx, advertiserId, size, (page-1)*size)
}

func (c *Campaign) Create(ctx ctx.Context, campaign *entity.Campaign, genText bool, prompt string) e.Error {
	_, err := c.advertiser.GetById(ctx, campaign.AdvertiserId)
	if err != nil {
		return err
	}

	if campaign.EndDate < campaign.StartDate {
		return badDateErr
	}

	curTime, err := c.time.Get(ctx)
	if err != nil {
		return e.InternalErr
	}

	if campaign.EndDate < curTime || campaign.StartDate < curTime {
		return badDateErr
	}

	if campaign.Targeting.AgeFrom != nil &&
		campaign.Targeting.AgeTo != nil &&
		*campaign.Targeting.AgeTo < *campaign.Targeting.AgeFrom {
		return badAgeErr
	}

	if genText {
		text, err := c.ai.GenText(prompt)
		if err != nil {
			return err
		}

		campaign.Text = text
	}

	err = c.moderate.Moderate(ctx, campaign.Text)
	if err != nil {
		return err
	}

	return c.campaign.Create(ctx, campaign)
}

func (c *Campaign) Update(ctx ctx.Context, campaign *entity.Campaign, genText bool, prompt string) e.Error {
	old, err := c.GetById(ctx, campaign)
	if err != nil {
		return err
	}

	curTime, err := c.time.Get(ctx)
	if err != nil {
		return e.InternalErr
	}

	isChanged := campaign.StartDate != old.StartDate ||
		campaign.EndDate != old.EndDate ||
		campaign.Billing.ImpressionsLimit != old.Billing.ImpressionsLimit ||
		campaign.Billing.ClicksLimit != old.Billing.ClicksLimit

	if !(curTime < old.StartDate) && isChanged {
		return badReqErr
	}

	if campaign.EndDate < curTime || campaign.StartDate < curTime {
		return badDateErr
	}

	if campaign.EndDate < campaign.StartDate {
		return badDateErr
	}

	if campaign.Targeting.AgeFrom != nil &&
		campaign.Targeting.AgeTo != nil &&
		*campaign.Targeting.AgeTo < *campaign.Targeting.AgeFrom {
		return badAgeErr
	}

	if genText {
		text, err := c.ai.GenText(prompt)
		if err != nil {
			return err
		}

		campaign.Text = text
	}

	err = c.moderate.Moderate(ctx, campaign.Text)
	if err != nil {
		return err
	}

	return c.campaign.Update(ctx, campaign)
}

func (c *Campaign) GetImage(ctx ctx.Context, campaign *entity.Campaign) (string, e.Error) {
	campaign, err := c.GetById(ctx, campaign)
	if err != nil {
		return "", err
	}

	if campaign.Image == "" {
		return "", e.New("This campaign hasn`t image.", e.NotFound)
	}

	return c.image.Get(ctx, campaign.Image)
}

func (c *Campaign) DownloadImage(ctx ctx.Context, campaign *entity.Campaign) ([]byte, e.Error) {
	campaign, err := c.GetById(ctx, campaign)
	if err != nil {
		return nil, err
	}

	if campaign.Image == "" {
		return nil, e.New("This campaign hasn`t image.", e.NotFound)
	}

	return c.image.Download(ctx, campaign.Image)
}

func (c *Campaign) UploadImage(ctx ctx.Context, campaign *entity.Campaign, image *entity.Image) (string, e.Error) {
	old, err := c.GetById(ctx, campaign)
	if err != nil {
		return "", err
	}

	if old.Image != "" {
		if err := c.image.Delete(ctx, old.Image); err != nil {
			return "", err
		}
	}

	name := image.Name
	parts := strings.Split(name, ".")

	if len(parts) < 2 {
		return "", e.New("Bad file name.", e.BadInput)
	}

	id := uuid.NewString()
	parts[0] = id

	name = strings.Join(parts, ".")
	image.Name = name

	url, err := c.image.Upload(ctx, image)
	if err != nil {
		return "", err
	}

	old.Image = name

	if err := c.campaign.Update(ctx, old); err != nil {
		return "", err
	}

	return url, nil
}

func (c *Campaign) DeleteImage(ctx ctx.Context, campaign *entity.Campaign) e.Error {
	campaign, err := c.GetById(ctx, campaign)
	if err != nil {
		return err
	}

	if campaign.Image == "" {
		return nil
	}

	if err := c.image.Delete(ctx, campaign.Image); err != nil {
		return err
	}

	campaign.Image = ""

	return c.campaign.Update(ctx, campaign)
}

func (c *Campaign) Delete(ctx ctx.Context, campaign *entity.Campaign) e.Error {
	campaign, err := c.GetById(ctx, campaign)
	if err != nil {
		return err
	}

	if campaign.Image != "" {
		if err := c.image.Delete(ctx, campaign.Image); err != nil {
			return err
		}
	}

	return c.campaign.Delete(ctx, campaign)
}
