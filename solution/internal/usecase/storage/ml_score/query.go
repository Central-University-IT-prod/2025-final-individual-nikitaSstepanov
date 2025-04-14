package mlscore

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func idQuery(clientId string, advertiserId string) (string, []interface{}) {
	builder := sq.Select("*").From(mlTable).
		Where(sq.Eq{"client_id": clientId, "advertiser_id": advertiserId})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()
	fmt.Println(args...)
	return query, args
}

func clientQuery(clientId string) (string, []interface{}) {
	builder := sq.Select("*").From(mlTable).
		Where(sq.Eq{"client_id": clientId})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()
	fmt.Println(args...)
	return query, args
}

func createQuery(score *entity.MlScore) (string, []interface{}) {
	builder := sq.Insert(mlTable).Columns("client_id", "advertiser_id", "score").
		Values(score.ClientId, score.AdvertiserId, score.Score)

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}

func updateQuery(score *entity.MlScore) (string, []interface{}) {
	builder := sq.Update(mlTable).
		Set("score", score.Score).
		Where(sq.Eq{"client_id": score.ClientId, "advertiser_id": score.AdvertiserId})

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	return query, args
}
