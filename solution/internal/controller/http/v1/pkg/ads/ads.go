package ads

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/httper"
	conv "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/converter"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Ads struct {
	usecase AdsUseCase
}

func New(uc AdsUseCase) *Ads {
	return &Ads{
		usecase: uc,
	}
}

func (a *Ads) Get(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Query("client_id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign, err := a.usecase.GetAd(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoAd(campaign)

	c.JSON(httper.StatusOK, result)
}

func (a *Ads) Click(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	var body dto.ClientId

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	click := &entity.Click{
		CampaignId: id,
		ClientId:   body.Id,
	}

	err := a.usecase.Click(ctx, click)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}
