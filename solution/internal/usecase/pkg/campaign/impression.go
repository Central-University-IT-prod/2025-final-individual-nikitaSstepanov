package campaign

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func (c *Campaign) GetAd(ctx ctx.Context, clientId string) (*entity.Campaign, e.Error) {
	client, err := c.client.GetById(ctx, clientId)
	if err != nil {
		return nil, err
	}

	day, err := c.time.Get(ctx)
	if err != nil {
		return nil, e.InternalErr
	}

	campaigns, err := c.campaign.GetAvailable(ctx, client, day)
	if err != nil {
		return nil, err
	}

	if len(campaigns) == 0 {
		return nil, e.New("Ad for user wasn`t found.", e.NotFound)
	}

	scores, err := c.score.GetClientScores(ctx, clientId)
	if err != nil {
		return nil, err
	}

	campaign, err := c.chooseCampaign(campaigns, scores)
	if err != nil {
		return nil, err
	}

	if campaign.Image != "" {
		url, err := c.GetImage(ctx, campaign)
		if err != nil {
			return nil, err
		}

		campaign.Image = url
	}

	_, err = c.campaign.GetImpression(ctx, campaign.Id, clientId)
	if err != nil && err.GetCode() != e.NotFound {
		return nil, err
	}

	if err == nil {
		return campaign, nil
	}

	impression := &entity.Impression{
		CampaignId: campaign.Id,
		ClientId:   clientId,
	}

	if err := c.campaign.CreateImpression(ctx, impression); err != nil {
		return nil, err
	}

	billing := campaign.Billing

	billing.ImpressionsCount += 1
	billing.SpentImpressions += billing.CostPerImpression

	if err := c.campaign.UpdateBilling(ctx, billing, campaign.Id); err != nil {
		return nil, err
	}

	daily, err := c.campaign.GetBillingByDay(ctx, campaign.Id, day)
	if err != nil && err.GetCode() != e.NotFound {
		return nil, err
	}

	if err != nil {
		daily = &entity.DailyBilling{
			Date:             day,
			ImpressionsCount: 1,
			SpentImpressions: billing.CostPerImpression,
		}

		err := c.campaign.CreateDailyBilling(ctx, daily, campaign.Id)
		if err != nil {
			return nil, err
		}
	} else {
		daily.ImpressionsCount += 1
		daily.SpentImpressions += billing.CostPerImpression
	}

	if err := c.campaign.UpdateDailyBilling(ctx, daily, campaign.Id); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (c *Campaign) chooseCampaign(campaigns []*entity.Campaign, scores map[string]int) (*entity.Campaign, e.Error) {
	result := &entity.Campaign{}
	cost := float32(0)
	score := 0

	for _, campaign := range campaigns {
		campScore, ok := scores[campaign.AdvertiserId]
		if !ok {
			campScore = 0
		}

		if campaign.Billing.CostPerImpression > cost {
			result = campaign
			cost = campaign.Billing.CostPerImpression
			score = campScore
		} else if campaign.Billing.CostPerImpression == cost && campScore > score {
			result = campaign
			cost = campaign.Billing.CostPerImpression
			score = campScore
		}
	}

	return result, nil
}
