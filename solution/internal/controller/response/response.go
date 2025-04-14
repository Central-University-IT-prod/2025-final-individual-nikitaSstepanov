package resp

import (
	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Message struct {
	Message string `json:"message"`
}

func NewMessage(msg string) *Message {
	return &Message{
		Message: msg,
	}
}

func AbortErrMsg(c *gin.Context, err e.Error) {
	ctx := ct.GetCtx(c)
	log := ctx.Logger()

	if err.GetCode() == e.Internal {
		log.Error("Something going wrong...", err.SlErr())
	} else {
		log.Info("Invalid input data", err.SlErr())
	}

	c.AbortWithStatusJSON(
		err.ToHttpCode(),
		err.ToJson(),
	)
}
