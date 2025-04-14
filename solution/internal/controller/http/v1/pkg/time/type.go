package time

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type TimeUseCase interface {
	Set(c ctx.Context, day int) e.Error
}
