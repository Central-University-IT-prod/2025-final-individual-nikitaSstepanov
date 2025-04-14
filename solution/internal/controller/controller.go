package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/httper"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/telegram"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase"
	"gopkg.in/telebot.v4"
)

type Controller struct {
	v1       *v1.Router
	telegram *telegram.Bot
	cfg      *Config
}

type Config struct {
	V1       v1.Config       `yaml:"v1"`
	Telegram telegram.Config `yaml:"telegram"`
	Mode     string          `yaml:"mode" env:"MODE" env-default:"DEBUG"`
}

func New(uc *usecase.UseCase, cfg *Config) *Controller {
	return &Controller{
		v1:       v1.New(uc, &cfg.V1),
		telegram: telegram.New(uc, &cfg.Telegram),
		cfg:      cfg,
	}
}

func (c *Controller) InitRoutes(ctx ctx.Context) *gin.Engine {
	setGinMode(c.cfg.Mode)

	router := gin.New()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(httper.StatusOK, "pong")
	})

	router.GET("/metrics", prometheusHandler())

	api := router.Group("/")
	{
		c.v1.InitRoutes(ctx, api)
	}

	return router
}

func (c *Controller) InitBot(ctx ctx.Context) *telebot.Bot {
	return c.telegram.InitBot(ctx)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func setGinMode(mode string) {
	switch mode {

	case "RELEASE":
		gin.SetMode(gin.ReleaseMode)

	case "TEST":
		gin.SetMode(gin.TestMode)

	case "DEBUG":
		gin.SetMode(gin.DebugMode)

	default:
		gin.SetMode(gin.DebugMode)

	}
}
