package campaign

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func (c *Campaign) Click(ctx ctx.Context, click *entity.Click) e.Error {
	campaign, err := c.campaign.GetById(ctx, click.CampaignId)
	if err != nil {
		return err
	}

	_, err = c.client.GetById(ctx, click.ClientId)
	if err != nil {
		return err
	}

	_, err = c.campaign.GetClick(ctx, click.CampaignId, click.ClientId)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	if err == nil {
		return nil
	}

	if err := c.campaign.CreateClick(ctx, click); err != nil {
		return err
	}

	billing := campaign.Billing

	billing.ClicksCount += 1
	billing.SpentClicks += billing.CostPerClick

	if err := c.campaign.UpdateBilling(ctx, billing, campaign.Id); err != nil {
		return err
	}

	day, err := c.time.Get(ctx)
	if err != nil {
		return e.InternalErr
	}

	daily, err := c.campaign.GetBillingByDay(ctx, campaign.Id, day)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	if err != nil {
		daily = &entity.DailyBilling{
			Date:        day,
			ClicksCount: 1,
			SpentClicks: billing.CostPerClick,
		}

		err := c.campaign.CreateDailyBilling(ctx, daily, campaign.Id)
		if err != nil {
			return err
		}
	} else {
		daily.ClicksCount += 1
		daily.SpentClicks += billing.CostPerClick
	}

	return c.campaign.UpdateDailyBilling(ctx, daily, campaign.Id)
}
