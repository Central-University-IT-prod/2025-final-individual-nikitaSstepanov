package entity

import "github.com/nikitaSstepanov/tools/client/pg"

type MlScore struct {
	ClientId     string
	AdvertiserId string
	Score        int
}

func (ml *MlScore) Scan(r pg.Row) error {
	return r.Scan(
		&ml.ClientId,
		&ml.AdvertiserId,
		&ml.Score,
	)
}
