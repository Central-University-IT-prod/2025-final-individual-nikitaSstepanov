package client

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func idQuery(id string) (string, []interface{}) {
	builder := sq.Select("*").From(clientsTable).Where(sq.Eq{"id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createQuery(clients []*entity.Client) (string, []interface{}) {
	builder := sq.Insert(clientsTable).Columns("id", "login", "age", "location", "gender")

	for _, client := range clients {
		builder = builder.Values(
			client.Id,
			client.Login,
			client.Age,
			client.Location,
			client.Gender,
		)
	}

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateQuery(client *entity.Client) (string, []interface{}) {
	builder := sq.Update(clientsTable).
		Set("login", client.Login).Set("age", client.Age).
		Set("location", client.Location).Set("gender", client.Gender).
		Where(sq.Eq{"id": client.Id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}
