package moderation

import (
	"strings"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type Moderation struct {
	blacklist BlackListStorage
	need      bool
}

func New(blacklist BlackListStorage, need bool) *Moderation {
	return &Moderation{
		blacklist: blacklist,
		need:      need,
	}
}

func (m *Moderation) Moderate(c ctx.Context, text string) e.Error {
	if m.need {
		list, err := m.blacklist.Get(c)
		if err != nil {
			return err
		}

		for _, word := range list {
			if strings.Contains(text, word) {
				return e.New("Invalid text.", e.BadInput)
			}
		}
	}

	return nil
}

func (m *Moderation) AddWord(c ctx.Context, word string) e.Error {
	return m.blacklist.Add(c, word)
}
