package moderation

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type BlackListStorage interface {
	Get(c ctx.Context) ([]string, e.Error)
	Add(c ctx.Context, word string) e.Error
}
