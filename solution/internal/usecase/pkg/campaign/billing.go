package campaign

import (
	"sort"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func (c *Campaign) CampaignBilling(ctx ctx.Context, id string) (*entity.Billing, e.Error) {
	return c.campaign.GetBilling(ctx, id)
}

func (c *Campaign) AdvertiserBilling(ctx ctx.Context, advertiserId string) (*entity.Billing, e.Error) {
	_, err := c.advertiser.GetById(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	campaigns, err := c.campaign.Get(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	result := &entity.Billing{}

	if len(campaigns) == 0 {
		return result, nil
	}

	for _, campaign := range campaigns {
		bill, err := c.campaign.GetBilling(ctx, campaign.Id)
		if err != nil {
			return nil, err
		}

		result.ImpressionsCount += bill.ImpressionsCount
		result.ClicksCount += bill.ClicksCount
		result.SpentImpressions += bill.SpentImpressions
		result.SpentClicks += bill.SpentClicks
	}

	return result, nil
}

func (c *Campaign) CampaignDailyBilling(ctx ctx.Context, id string) ([]*entity.DailyBilling, e.Error) {
	return c.campaign.GetDailyBilling(ctx, id)
}

func (c *Campaign) CampaignDailyBillingWithPagination(ctx ctx.Context, id string, size, page int) ([]*entity.DailyBilling, e.Error) {
	return c.campaign.GetDailyBillingWithPagination(ctx, id, size, (page-1)*size)
}

func (c *Campaign) AdvertiserDailyBill(ctx ctx.Context, advertiserId string) ([]*entity.DailyBilling, e.Error) {
	_, err := c.advertiser.GetById(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	campaigns, err := c.campaign.Get(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	bills := make(map[int]*entity.DailyBilling)

	for _, campaign := range campaigns {
		dailyBills, err := c.campaign.GetDailyBilling(ctx, campaign.Id)
		if err != nil {
			return nil, err
		}

		for _, dailyBill := range dailyBills {
			bill, ok := bills[dailyBill.Date]
			if ok {
				bill.ImpressionsCount += dailyBill.ImpressionsCount
				bill.ClicksCount += dailyBill.ClicksCount
				bill.SpentImpressions += dailyBill.SpentImpressions
				bill.SpentClicks += dailyBill.SpentClicks
			} else {
				toAdd := &entity.DailyBilling{
					Date:             dailyBill.Date,
					ImpressionsCount: dailyBill.ImpressionsCount,
					ClicksCount:      dailyBill.ClicksCount,
					SpentImpressions: dailyBill.SpentImpressions,
					SpentClicks:      dailyBill.SpentClicks,
				}

				bills[dailyBill.Date] = toAdd
			}
		}
	}

	result := make([]*entity.DailyBilling, 0)

	for _, bill := range bills {
		result = append(result, bill)
	}

	return result, nil
}

func (c *Campaign) AdvertiserDailyBillWithPagination(ctx ctx.Context, advertiserId string, size, page int) ([]*entity.DailyBilling, e.Error) {
	_, err := c.advertiser.GetById(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	campaigns, err := c.campaign.Get(ctx, advertiserId)
	if err != nil {
		return nil, err
	}

	bills := make(map[int]*entity.DailyBilling)

	for _, campaign := range campaigns {
		dailyBills, err := c.campaign.GetDailyBilling(ctx, campaign.Id)
		if err != nil {
			return nil, err
		}

		for _, dailyBill := range dailyBills {
			bill, ok := bills[dailyBill.Date]
			if ok {
				bill.ImpressionsCount += dailyBill.ImpressionsCount
				bill.ClicksCount += dailyBill.ClicksCount
				bill.SpentImpressions += dailyBill.SpentImpressions
				bill.SpentClicks += dailyBill.SpentClicks
			} else {
				toAdd := &entity.DailyBilling{
					Date:             dailyBill.Date,
					ImpressionsCount: dailyBill.ImpressionsCount,
					ClicksCount:      dailyBill.ClicksCount,
					SpentImpressions: dailyBill.SpentImpressions,
					SpentClicks:      dailyBill.SpentClicks,
				}

				bills[dailyBill.Date] = toAdd
			}
		}
	}

	result := make([]*entity.DailyBilling, 0)

	for _, bill := range bills {
		result = append(result, bill)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	if (page-1)*size+size < len(result) {
		return result[(page-1)*size : (page-1)*size+size], nil
	} else {
		return result[(page-1)*size:], nil
	}
}
