package campaign

import (
	"fmt"

	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"

	sq "github.com/Masterminds/squirrel"
)

func availableQuery(client *entity.Client, day int) (string, []interface{}) {
	builder := sq.Select("*").From(fmt.Sprintf("%s AS c", campaignTable)).
		LeftJoin(fmt.Sprintf("%s AS t ON c.id = t.campaign_id", targetTable)).
		Join(fmt.Sprintf("%s AS b ON c.id = b.campaign_id", billingTable)).
		Where(
			sq.And{
				sq.LtOrEq{
					"start_date": day,
				},
				sq.GtOrEq{
					"end_date": day,
				},
				sq.Or{
					sq.Eq{
						"t.gender": nil,
					},
					sq.Eq{
						"t.gender": client.Gender,
					},
					sq.Eq{
						"t.gender": types.ALL,
					},
				},
				sq.Or{
					sq.Eq{
						"t.age_from": nil,
					},
					sq.LtOrEq{
						"t.age_from": client.Age,
					},
				},
				sq.Or{
					sq.Eq{
						"t.age_to": nil,
					},
					sq.GtOrEq{
						"t.age_to": client.Age,
					},
				},
				sq.Or{
					sq.Eq{
						"t.location": nil,
					},
					sq.Eq{
						"t.location": client.Location,
					},
				},
			},
		)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()
	query += " AND (b.impressions_count < b.impressions_limit)"
	query += " AND (b.clicks_count < b.clicks_limit)"
	fmt.Println(query)
	return query, args
}

func idQuery(id string) (string, []interface{}) {
	builder := sq.Select("*").From(fmt.Sprintf("%s AS c", campaignTable)).
		LeftJoin(fmt.Sprintf("%s AS t ON c.id = t.campaign_id", targetTable)).
		Join(fmt.Sprintf("%s AS b ON c.id = b.campaign_id", billingTable)).
		Where(sq.Eq{"id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func getQuery(advertiserId string) (string, []interface{}) {
	builder := sq.Select("*").From(fmt.Sprintf("%s AS c", campaignTable)).
		LeftJoin(fmt.Sprintf("%s AS t ON c.id = t.campaign_id", targetTable)).
		Join(fmt.Sprintf("%s AS b ON c.id = b.campaign_id", billingTable)).
		Where(sq.Eq{"advertiser_id": advertiserId})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func paginationQuery(advertiserId string, limit int, offset int) (string, []interface{}) {
	builder := sq.Select("*").From(fmt.Sprintf("%s AS c", campaignTable)).
		LeftJoin(fmt.Sprintf("%s AS t ON c.id = t.campaign_id", targetTable)).
		Join(fmt.Sprintf("%s AS b ON c.id = b.campaign_id", billingTable)).
		Where(sq.Eq{"advertiser_id": advertiserId}).Limit(uint64(limit)).Offset(uint64(offset))

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func billQuery(id string) (string, []interface{}) {
	builder := sq.Select("*").From(billingTable).
		Where(sq.Eq{"campaign_id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func dailyBillQuery(id string) (string, []interface{}) {
	builder := sq.Select("*").From(dailyBillingTable).
		Where(sq.Eq{"campaign_id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func dailyPaginationQuery(id string, limit, offset int) (string, []interface{}) {
	builder := sq.Select("*").From(dailyBillingTable).
		Where(sq.Eq{"campaign_id": id}).OrderBy("date").Limit(uint64(limit)).Offset(uint64(offset))

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func dailyQuery(id string, day int) (string, []interface{}) {
	builder := sq.Select("*").From(dailyBillingTable).
		Where(sq.Eq{"campaign_id": id, "date": day})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createDailyQuery(billing *entity.DailyBilling, id string) (string, []interface{}) {
	builder := sq.Insert(dailyBillingTable).
		Columns("campaign_id", "date").
		Values(id, billing.Date)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateDailyQuery(billing *entity.DailyBilling, id string) (string, []interface{}) {
	builder := sq.Update(dailyBillingTable).
		Set("impressions_count", billing.ImpressionsCount).
		Set("clicks_count", billing.ClicksCount).
		Set("spent_impressions", billing.SpentImpressions).
		Set("spent_clicks", billing.SpentClicks).
		Where(sq.Eq{"campaign_id": id, "date": billing.Date})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func impressionQuery(campaignId, clientId string) (string, []interface{}) {
	builder := sq.Select("*").From(impressionsTable).
		Where(sq.Eq{"campaign_id": campaignId, "client_id": clientId})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func clickQuery(campaignId, clientId string) (string, []interface{}) {
	builder := sq.Select("*").From(clicksTable).
		Where(sq.Eq{"campaign_id": campaignId, "client_id": clientId})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createQuery(campaign *entity.Campaign) (string, []interface{}) {
	builder := sq.Insert(campaignTable).
		Columns("title", "text", "start_date", "end_date", "advertiser_id", "image").
		Values(campaign.Title, campaign.Text, campaign.StartDate, campaign.EndDate, campaign.AdvertiserId, campaign.Image).
		Suffix("RETURNING id")

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createTargetingQuery(targeting *entity.Targeting, id string) (string, []interface{}) {
	builder := sq.Insert(targetTable).
		Columns("campaign_id", "gender", "age_from", "age_to", "location").
		Values(id, targeting.Gender, targeting.AgeFrom, targeting.AgeTo, targeting.Location)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createBillingQuery(billing *entity.Billing, id string) (string, []interface{}) {
	builder := sq.Insert(billingTable).
		Columns("campaign_id", "impressions_limit", "clicks_limit", "cost_per_impression", "cost_per_click").
		Values(id, billing.ImpressionsLimit, billing.ClicksLimit, billing.CostPerImpression, billing.CostPerClick)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createImpressionQuery(impression *entity.Impression) (string, []interface{}) {
	builder := sq.Insert(impressionsTable).
		Columns("campaign_id", "client_id").
		Values(impression.CampaignId, impression.ClientId)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createClickQuery(click *entity.Click) (string, []interface{}) {
	builder := sq.Insert(clicksTable).
		Columns("campaign_id", "client_id").
		Values(click.CampaignId, click.ClientId)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateQuery(campaign *entity.Campaign) (string, []interface{}) {
	builder := sq.Update(campaignTable).
		Set("title", campaign.Title).
		Set("text", campaign.Text).
		Set("start_date", campaign.StartDate).
		Set("end_date", campaign.EndDate).
		Set("image", campaign.Image).
		Where(sq.Eq{"id": campaign.Id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateBillingQuery(billing *entity.Billing, id string) (string, []interface{}) {
	builder := sq.Update(billingTable).
		Set("impressions_limit", billing.ImpressionsLimit).
		Set("clicks_limit", billing.ClicksLimit).
		Set("cost_per_impression", billing.CostPerImpression).
		Set("cost_per_click", billing.CostPerClick).
		Where(sq.Eq{"campaign_id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateCountsQuery(billing *entity.Billing, id string) (string, []interface{}) {
	builder := sq.Update(billingTable).
		Set("impressions_count", billing.ImpressionsCount).
		Set("clicks_count", billing.ClicksCount).
		Set("spent_impressions", billing.SpentImpressions).
		Set("spent_clicks", billing.SpentClicks).
		Where(sq.Eq{"campaign_id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateTargetingQuery(targeting *entity.Targeting, id string) (string, []interface{}) {
	builder := sq.Update(targetTable).
		Set("gender", targeting.Gender).
		Set("age_from", targeting.AgeFrom).
		Set("age_to", targeting.AgeTo).
		Set("location", targeting.Location).
		Where(sq.Eq{"campaign_id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func deleteQuery(campaign *entity.Campaign) (string, []interface{}) {
	builder := sq.Delete(campaignTable).Where(sq.Eq{"id": campaign.Id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}
