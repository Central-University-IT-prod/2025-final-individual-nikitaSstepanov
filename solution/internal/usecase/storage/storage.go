package storage

import (
	"github.com/nikitaSstepanov/tools"
	"github.com/nikitaSstepanov/tools/client/pg"
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/advertiser"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/blacklist"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/campaign"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/client"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/image"
	mlscore "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/ml_score"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/telegram"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage/time"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/client/minio"
)

type Storage struct {
	Advertiser *advertiser.Advertiser
	Campaign   *campaign.Campaign
	Client     *client.Client
	Score      *mlscore.MlScore
	Image      *image.Image
	Telegram   *telegram.Telegram
	Blacklist  *blacklist.Blacklist
	Time       *time.Time
	pg         pg.Client
	rs         rs.Client
	mn         minio.Client
}

type Config struct {
	Minio minio.Config `yaml:"minio"`
}

func New(c ctx.Context, cfg *Config) *Storage {
	postgres := connectPg(c)
	redis := connectRs(c)
	minio := connectMinio(c, &cfg.Minio)

	return &Storage{
		Advertiser: advertiser.New(postgres),
		Campaign:   campaign.New(postgres),
		Client:     client.New(postgres),
		Score:      mlscore.New(postgres),
		Time:       time.New(c, redis),
		Image:      image.New(minio, cfg.Minio.Bucket),
		Telegram:   telegram.New(redis),
		Blacklist:  blacklist.New(redis),
		pg:         postgres,
		rs:         redis,
		mn:         minio,
	}
}

func (s *Storage) Close() e.Error {
	s.pg.Close()
	return e.E(s.rs.Close())
}

func connectPg(c ctx.Context) pg.Client {
	log := c.Logger()

	postgres, err := tools.Pg()
	if err != nil {
		log.Error("Can`t connect to postgres.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Postgres is connected.")
	}

	if err := postgres.RegisterTypes(pgTypes); err != nil {
		log.Error("Can`t register custom types.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Custom types are registered.")
	}

	return postgres
}

func connectRs(c ctx.Context) rs.Client {
	log := c.Logger()

	redis, err := tools.Redis()
	if err != nil {
		log.Error("Can`t connect to redis.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Redis is connected.")
	}

	return redis
}

func connectMinio(c ctx.Context, cfg *minio.Config) minio.Client {
	log := c.Logger()

	minio, err := minio.New(c, cfg)
	if err != nil {
		log.Error("Can`t connect to minio.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Minio is connected.")
	}

	return minio
}
