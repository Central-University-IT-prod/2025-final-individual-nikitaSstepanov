package admin

import (
	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Admin struct {
	moderation ModerationUseCase
}

func New(moderation ModerationUseCase) *Admin {
	return &Admin{
		moderation: moderation,
	}
}

func (a *Admin) UpdateBlacklist(c *gin.Context) {
	ctx := ct.GetCtx(c)

	word := c.Query("word")

	if word == "" {
		resp.AbortErrMsg(c, e.New("Word can`t be empty", e.BadInput))
		return
	}

	err := a.moderation.AddWord(ctx, word)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(204, nil)
}
