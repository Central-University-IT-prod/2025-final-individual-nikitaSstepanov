package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/middleware"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/admin"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/ads"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/advertiser"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/campaign"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/client"
	mlscore "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/ml_score"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/stats"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/pkg/time"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/usecase"
)

type Router struct {
	advertiser AdvertiserHandler
	ads        AdsHandler
	campaign   CampaignHandler
	client     ClientHandler
	score      ScoreHandler
	stats      StatsHandler
	time       TimeHandler
	admin      AdminHandler
	mid        Middleware
}

type Config struct{}

func New(uc *usecase.UseCase, cfg *Config) *Router {
	return &Router{
		advertiser: advertiser.New(uc.Advertiser),
		ads:        ads.New(uc.Campaign),
		campaign:   campaign.New(uc.Campaign),
		client:     client.New(uc.Client),
		time:       time.New(uc.Time),
		score:      mlscore.New(uc.Score),
		stats:      stats.New(uc.Campaign),
		mid:        middleware.New(),
		admin:      admin.New(uc.Blacklist),
	}
}

func (r *Router) InitRoutes(ctx ctx.Context, h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/")
	{
		router.Use(r.mid.InitLogger(ctx))

		r.initAdvertiserRoutes(router)
		r.initAdsRoutes(router)
		r.initClientRoutes(router)
		r.initTimeRoutes(router)
		r.initStatsRoutes(router)
		r.initScoreRoutes(router)
		router.POST("/blacklist", r.admin.UpdateBlacklist)
	}

	return router
}

func (r *Router) initAdvertiserRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/advertisers")
	{
		router.GET("/:id", r.advertiser.GetById)
		router.POST("/bulk", r.advertiser.Bulk)

		id := router.Group("/:id")
		{
			id.POST("/campaigns", r.campaign.Create)
			id.GET("/campaigns", r.campaign.Get)
			id.GET("/campaigns/:camId", r.campaign.GetById)
			id.PUT("/campaigns/:camId", r.campaign.Update)
			id.GET("/campaigns/:camId/image", r.campaign.GetImage)
			id.POST("/campaigns/:camId/image/upload", r.campaign.UploadImage)
			id.DELETE("/campaigns/:camId/image/delete", r.campaign.DeleteImage)
			id.DELETE("/campaigns/:camId", r.campaign.Delete)
		}
	}

	return router
}

func (r *Router) initAdsRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	h.GET("/ads", r.ads.Get)
	h.POST("/ads/:id/click", r.ads.Click)

	return h
}

func (r *Router) initClientRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/clients")
	{
		router.GET("/:id", r.client.GetById)
		router.POST("/bulk", r.client.Bulk)
	}

	return router
}

func (r *Router) initStatsRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/stats")
	{
		router.GET("/campaigns/:id", r.stats.Campaign)
		router.GET("/advertisers/:id/campaigns", r.stats.Advertiser)
		router.GET("/campaigns/:id/daily", r.stats.CampaignDaily)
		router.GET("/advertisers/:id/campaigns/daily", r.stats.AdvertiserDaily)
	}

	return router
}

func (r *Router) initScoreRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	h.POST("/ml-scores", r.score.Manage)

	return h
}

func (r *Router) initTimeRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	h.POST("/time/advance", r.time.Advance)

	return h
}
