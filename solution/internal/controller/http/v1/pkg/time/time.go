package time

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/httper"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Time struct {
	usecase TimeUseCase
}

func New(uc TimeUseCase) *Time {
	return &Time{
		usecase: uc,
	}
}

func (t *Time) Advance(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body dto.AdvanceTime

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	err := t.usecase.Set(ctx, body.Day)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusOK, body)
}
