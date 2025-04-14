package admin

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type ModerationUseCase interface {
	AddWord(c ctx.Context, word string) e.Error
}
