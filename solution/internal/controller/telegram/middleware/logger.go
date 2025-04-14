package middleware

import (
	"time"

	goctx "github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/sl"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
	"gopkg.in/telebot.v4"
)

func (m *Middleware) InitLogger(c goctx.Context) telebot.MiddlewareFunc {
	log := c.Logger()

	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			c := goctx.New(log)

			ctx.Set(ct.CtxKey, c)

			t1 := time.Now()

			err := next(ctx)
			if err != nil {
				return err
			}

			entry := log.With(

				sl.Int64Attr("tg_id", ctx.Sender().ID),
				sl.StringAttr("duration", time.Since(t1).String()),
			)

			if len(ctx.Message().Entities) != 0 {
				entry = entry.With(
					sl.AnyAttr("type", ctx.Message().Entities[0].Type),
					sl.StringAttr("text", ctx.Text()),
				)
			}

			entry.Info("telegram request completed")

			return nil
		}
	}
}
