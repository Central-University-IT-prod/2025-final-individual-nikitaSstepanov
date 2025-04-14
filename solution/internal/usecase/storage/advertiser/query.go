package advertiser

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func idQuery(id string) (string, []interface{}) {
	builder := sq.Select("*").From(advertisersTable).Where(sq.Eq{"id": id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func createQuery(advertisers []*entity.Advertiser) (string, []interface{}) {
	builder := sq.Insert(advertisersTable).Columns("id", "name")

	for _, advertiser := range advertisers {
		builder = builder.Values(
			advertiser.Id,
			advertiser.Name,
		)
	}

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateQuery(advertiser *entity.Advertiser) (string, []interface{}) {
	builder := sq.Update(advertisersTable).
		Set("name", advertiser.Name).
		Where(sq.Eq{"id": advertiser.Id})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}
