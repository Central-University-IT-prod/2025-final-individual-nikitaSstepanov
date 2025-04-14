package entity

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"
)

type Campaign struct {
	Id           string
	AdvertiserId string
	Title        string
	Text         string
	StartDate    int
	EndDate      int
	Image        string
	Targeting    *Targeting
	Billing      *Billing
}

type Targeting struct {
	Gender   *types.Gender
	AgeFrom  *int
	AgeTo    *int
	Location *string
}

type Billing struct {
	ImpressionsLimit  int
	ClicksLimit       int
	CostPerImpression float32
	CostPerClick      float32
	ImpressionsCount  int
	ClicksCount       int
	SpentImpressions  float32
	SpentClicks       float32
}

type DailyBilling struct {
	Date             int
	ImpressionsCount int
	ClicksCount      int
	SpentImpressions float32
	SpentClicks      float32
}

type Impression struct {
	CampaignId string
	ClientId   string
}

type Click struct {
	CampaignId string
	ClientId   string
}

func (c *Campaign) Scan(r pg.Row) error {
	var t Targeting
	var b Billing
	var id *string

	err := r.Scan(
		&c.Id,
		&c.Title,
		&c.Text,
		&c.StartDate,
		&c.EndDate,
		&c.AdvertiserId,
		&c.Image,
		&id,
		&t.Gender,
		&t.AgeFrom,
		&t.AgeTo,
		&t.Location,
		&id,
		&b.ImpressionsLimit,
		&b.ClicksLimit,
		&b.CostPerImpression,
		&b.CostPerClick,
		&b.ImpressionsCount,
		&b.ClicksCount,
		&b.SpentImpressions,
		&b.SpentClicks,
	)

	if err != nil {
		return err
	}

	c.Targeting = &t
	c.Billing = &b

	return nil
}

func (db *DailyBilling) Scan(r pg.Row) error {
	var id string

	return r.Scan(
		&id,
		&db.Date,
		&db.ImpressionsCount,
		&db.ClicksCount,
		&db.SpentImpressions,
		&db.SpentClicks,
	)
}

func (b *Billing) Scan(r pg.Row) error {
	var id string

	return r.Scan(
		&id,
		&b.ImpressionsLimit,
		&b.ClicksLimit,
		&b.CostPerImpression,
		&b.CostPerClick,
		&b.ImpressionsCount,
		&b.ClicksCount,
		&b.SpentImpressions,
		&b.SpentClicks,
	)
}

func (i *Impression) Scan(r pg.Row) error {
	return r.Scan(
		&i.ClientId,
		&i.CampaignId,
	)
}

func (c *Click) Scan(r pg.Row) error {
	return r.Scan(
		&c.ClientId,
		&c.CampaignId,
	)
}
