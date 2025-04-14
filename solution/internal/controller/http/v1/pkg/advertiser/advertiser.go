package advertiser

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/httper"
	conv "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/converter"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Advertiser struct {
	usecase AdvertiserUseCase
}

func New(uc AdvertiserUseCase) *Advertiser {
	return &Advertiser{
		usecase: uc,
	}
}

func (a *Advertiser) GetById(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	advertiser, err := a.usecase.GetById(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoAdvertiser(advertiser)

	c.JSON(httper.StatusOK, result)
}

func (a *Advertiser) Bulk(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body []dto.Advertiser

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	for _, advertiser := range body {
		if err := validator.Struct(advertiser); err != nil {
			resp.AbortErrMsg(c, err)
			return
		}
	}

	advertisers := conv.BulkAdvertiser(body)

	if err := a.usecase.Bulk(ctx, advertisers); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusCreated, body)
}
