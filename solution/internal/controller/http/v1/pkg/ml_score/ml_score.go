package mlscore

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/httper"
	conv "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/converter"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type MlScore struct {
	usecase ScoreUseCase
}

func New(uc ScoreUseCase) *MlScore {
	return &MlScore{
		usecase: uc,
	}
}

func (ml *MlScore) Manage(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body dto.MlScore

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	score := conv.EntityScore(body)

	err := ml.usecase.Manage(ctx, score)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusOK, okMsg)
}
