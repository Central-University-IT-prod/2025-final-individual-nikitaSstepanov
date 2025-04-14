package blacklist

import (
	"context"
	"encoding/json"

	"github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type Blacklist struct {
	redis redis.Client
}

type List struct {
	Words []string `redis:"blacklist"`
}

func (c *List) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *List) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func New(redis redis.Client) *Blacklist {
	list := List{
		Words: []string{},
	}

	err := redis.Set(context.Background(), "blacklist", &list, 0).Err()
	if err != nil {
		panic(err)
	}

	return &Blacklist{
		redis: redis,
	}
}

func (b *Blacklist) Get(c ctx.Context) ([]string, e.Error) {
	var list List

	err := b.redis.Get(c, "blacklist").Scan(&list)
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	return list.Words, nil
}

func (b *Blacklist) Add(c ctx.Context, word string) e.Error {
	list, err := b.Get(c)
	if err != nil {
		return e.InternalErr.WithErr(err)
	}

	list = append(list, word)

	new := List{
		Words: list,
	}

	setErr := b.redis.Set(c, "blacklist", &new, 0).Err()
	if setErr != nil {
		return e.InternalErr.WithErr(setErr)
	}

	return nil
}
