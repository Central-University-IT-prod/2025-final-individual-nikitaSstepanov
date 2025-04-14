package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/httper"
	conv "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/converter"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Stats struct {
	usecase BillingUseCase
}

func New(uc BillingUseCase) *Stats {
	return &Stats{
		usecase: uc,
	}
}

func (s *Stats) Campaign(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	billing, err := s.usecase.CampaignBilling(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoStats(billing)

	c.JSON(httper.StatusOK, result)
}

func (s *Stats) Advertiser(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	bill, err := s.usecase.AdvertiserBilling(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoStats(bill)

	c.JSON(httper.StatusOK, result)
}

func (s *Stats) CampaignDaily(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	bills, err := s.usecase.CampaignDailyBilling(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := make([]*dto.DailyStats, 0)

	for _, bill := range bills {
		toAdd := conv.DtoDailyStats(bill)

		result = append(result, toAdd)
	}

	c.JSON(httper.StatusOK, result)
}

func (s *Stats) AdvertiserDaily(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	bills, err := s.usecase.AdvertiserDailyBill(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := make([]*dto.DailyStats, 0)

	for _, bill := range bills {
		toAdd := conv.DtoDailyStats(bill)

		result = append(result, toAdd)
	}

	c.JSON(httper.StatusOK, result)
}
