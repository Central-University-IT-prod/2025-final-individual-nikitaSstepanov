package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
)

type Middleware interface {
	InitLogger(c ctx.Context) gin.HandlerFunc
}

type AdvertiserHandler interface {
	GetById(c *gin.Context)
	Bulk(c *gin.Context)
}

type AdminHandler interface {
	UpdateBlacklist(c *gin.Context)
}

type AdsHandler interface {
	Get(c *gin.Context)
	Click(c *gin.Context)
}

type CampaignHandler interface {
	GetById(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	GetImage(c *gin.Context)
	UploadImage(c *gin.Context)
	DeleteImage(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type StatsHandler interface {
	Campaign(c *gin.Context)
	Advertiser(c *gin.Context)
	CampaignDaily(c *gin.Context)
	AdvertiserDaily(c *gin.Context)
}

type ClientHandler interface {
	GetById(c *gin.Context)
	Bulk(c *gin.Context)
}

type ScoreHandler interface {
	Manage(c *gin.Context)
}

type TimeHandler interface {
	Advance(c *gin.Context)
}
