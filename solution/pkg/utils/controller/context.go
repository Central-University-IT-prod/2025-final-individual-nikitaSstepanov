package ct

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/sl"
	"gopkg.in/telebot.v4"
)

const (
	CtxKey = "ctx"
)

func GetCtx(c *gin.Context) ctx.Context {
	if c, ok := c.Get(CtxKey); ok {
		return c.(ctx.Context)
	}

	return ctx.New(sl.Default())
}

func GetCtxTg(c telebot.Context) ctx.Context {
	return c.Get(CtxKey).(ctx.Context)
}

func GetL(c *gin.Context) *sl.Logger {
	ctx := GetCtx(c)

	return ctx.Logger()
}
