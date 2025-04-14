package converter

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func DtoStats(billing *entity.Billing) *dto.Stats {
	var conversion float32

	if billing.ImpressionsCount == 0 {
		conversion = 0
	} else if billing.ClicksCount == 0 {
		conversion = 0
	} else {
		conversion = float32(billing.ClicksCount) / float32(billing.ImpressionsCount) * 100
	}

	return &dto.Stats{
		ImpressionsCount: billing.ImpressionsCount,
		ClicksCount:      billing.ClicksCount,
		Conversion:       conversion,
		SpentImpressions: billing.SpentImpressions,
		SpentClicks:      billing.SpentClicks,
		SpentTotal:       billing.SpentClicks + billing.SpentImpressions,
	}
}

func DtoDailyStats(billing *entity.DailyBilling) *dto.DailyStats {
	var conversion float32

	if billing.ImpressionsCount == 0 {
		conversion = 0
	} else if billing.ClicksCount == 0 {
		conversion = 0
	} else {
		conversion = float32(billing.ClicksCount) / float32(billing.ImpressionsCount) * 100
	}

	return &dto.DailyStats{
		Date:             billing.Date,
		ImpressionsCount: billing.ImpressionsCount,
		ClicksCount:      billing.ClicksCount,
		Conversion:       conversion,
		SpentImpressions: billing.SpentImpressions,
		SpentClicks:      billing.SpentClicks,
		SpentTotal:       billing.SpentClicks + billing.SpentImpressions,
	}
}
