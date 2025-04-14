package time

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type Time struct {
	storage TimeStorage
}

func New(storage TimeStorage) *Time {
	return &Time{
		storage: storage,
	}
}

func (t *Time) Get(c ctx.Context) (int, e.Error) {
	return t.storage.Get(c)
}

func (t *Time) Set(c ctx.Context, day int) e.Error {
	return t.storage.Set(c, day)
}
