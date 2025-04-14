package entity

import "github.com/nikitaSstepanov/tools/client/pg"

type Advertiser struct {
	Id   string
	Name string
}

func (a *Advertiser) Scan(r pg.Row) error {
	return r.Scan(
		&a.Id,
		&a.Name,
	)
}
