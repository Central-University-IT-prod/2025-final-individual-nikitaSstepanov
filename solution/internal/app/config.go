package app

import (
	config "github.com/nikitaSstepanov/tools/configurator"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase/storage"
)

type Config struct {
	Controller controller.Config `yaml:"controller"`
	UseCase    usecase.Config    `yaml:"usecase"`
	Storage    storage.Config    `yaml:"storage"`
}

func getConfig() (*Config, error) {
	var cfg Config

	if err := config.Get(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
