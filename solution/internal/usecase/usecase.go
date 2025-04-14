package usecase

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/ai"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/moderation"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/pkg/advertiser"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/pkg/campaign"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/pkg/client"
	mlscore "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/pkg/ml_score"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/pkg/telegram"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/pkg/time"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage"
)

type UseCase struct {
	Advertiser *advertiser.Advertiser
	Campaign   *campaign.Campaign
	Client     *client.Client
	Score      *mlscore.MlScore
	Time       *time.Time
	Telegram   *telegram.Telegram
	Blacklist  *moderation.Moderation
	Ai         *ai.Ai
}

type Config struct {
	Ai             ai.Config `json:"ai"`
	NeedModeration bool      `env:"NEED_MODERATION" env-default:"false"`
}

func New(store *storage.Storage, cfg *Config) *UseCase {
	time := time.New(store.Time)
	score := mlscore.New(store.Score)
	ai := ai.New(&cfg.Ai)
	moder := moderation.New(store.Blacklist, cfg.NeedModeration)

	return &UseCase{
		Time:       time,
		Client:     client.New(store.Client),
		Score:      score,
		Advertiser: advertiser.New(store.Advertiser),
		Campaign:   campaign.New(store.Campaign, store.Advertiser, time, store.Client, store.Image, score, ai, moder),
		Telegram:   telegram.New(store.Telegram, store.Campaign),
		Ai:         ai,
		Blacklist:  moder,
	}
}
