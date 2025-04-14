package advertiser

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Advertiser struct {
	storage AdvertiserStorage
}

func New(storage AdvertiserStorage) *Advertiser {
	return &Advertiser{
		storage: storage,
	}
}

func (a *Advertiser) GetById(c ctx.Context, id string) (*entity.Advertiser, e.Error) {
	return a.storage.GetById(c, id)
}

func (a *Advertiser) Bulk(c ctx.Context, advertisers []*entity.Advertiser) e.Error {
	toCreate := make([]*entity.Advertiser, 0)

	for _, advertiser := range advertisers {
		_, err := a.storage.GetById(c, advertiser.Id)
		if err != nil && err.GetCode() != e.NotFound {
			return err
		}

		if err != nil && err.GetCode() == e.NotFound {
			toCreate = append(toCreate, advertiser)
		} else {
			err := a.storage.Update(c, advertiser)
			if err != nil {
				return err
			}
		}
	}

	if len(toCreate) != 0 {
		return a.storage.Create(c, toCreate)
	}

	return nil
}
