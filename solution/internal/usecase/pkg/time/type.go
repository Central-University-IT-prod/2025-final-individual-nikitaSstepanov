package time

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type TimeStorage interface {
	Get(c ctx.Context) (int, e.Error)
	Set(c ctx.Context, day int) e.Error
}
