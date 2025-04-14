package entity

import (
	"github.com/nikitaSstepanov/tools/client/pg"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity/types"
)

type Client struct {
	Id       string
	Login    string
	Age      int
	Location string
	Gender   types.Gender
}

func (c *Client) Scan(r pg.Row) error {
	return r.Scan(
		&c.Id,
		&c.Login,
		&c.Age,
		&c.Location,
		&c.Gender,
	)
}
