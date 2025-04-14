package client

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

type Client struct {
	storage ClientStorage
}

func New(storage ClientStorage) *Client {
	return &Client{
		storage: storage,
	}
}

func (cl *Client) GetById(c ctx.Context, id string) (*entity.Client, e.Error) {
	return cl.storage.GetById(c, id)
}

func (cl *Client) Bulk(c ctx.Context, clients []*entity.Client) e.Error {
	toCreate := make([]*entity.Client, 0)

	for _, client := range clients {
		_, err := cl.storage.GetById(c, client.Id)
		if err != nil && err.GetCode() != e.NotFound {
			return err
		}

		if err != nil && err.GetCode() == e.NotFound {
			toCreate = append(toCreate, client)
		} else {
			err := cl.storage.Update(c, client)
			if err != nil {
				return err
			}
		}
	}

	if len(toCreate) != 0 {
		return cl.storage.Create(c, toCreate)
	}

	return nil
}
