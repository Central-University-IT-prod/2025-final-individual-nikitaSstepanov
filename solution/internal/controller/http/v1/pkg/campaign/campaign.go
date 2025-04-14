package campaign

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	conv "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/converter"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Campaign struct {
	usecase CampaignUseCase
}

func New(uc CampaignUseCase) *Campaign {
	return &Campaign{
		usecase: uc,
	}
}

func (ca *Campaign) GetById(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	camId := c.Param("camId")

	if err := validator.UUID(camId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	campaign, err := ca.usecase.GetById(ctx, campaign)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoCampaign(campaign)

	c.JSON(httper.StatusOK, result)
}

func (ca *Campaign) Get(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "5"))
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	campaigns, getErr := ca.usecase.Get(ctx, id, size, page)
	if getErr != nil {
		resp.AbortErrMsg(c, getErr)
		return
	}

	result := make([]map[string]interface{}, 0)

	for _, campaign := range campaigns {
		toAdd := conv.DtoCampaign(campaign)

		result = append(result, toAdd)
	}

	c.JSON(httper.StatusOK, result)
}

func (ca *Campaign) Create(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	var body dto.CreateCampaign

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign := conv.CreateCampaign(body, id)

	needGen := false
	prompt := ""

	if body.GenText != nil {
		needGen = *body.GenText

		if body.Prompt == nil {
			resp.AbortErrMsg(c, e.New("If you use text generating, prompt is required", e.BadInput))
			return
		}
		prompt = *body.Prompt
	}
	if needGen {
		campaign.Text = prompt
	}

	err := ca.usecase.Create(ctx, campaign, needGen, prompt)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoCampaign(campaign)

	c.JSON(201, result)
}

func (ca *Campaign) Update(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	camId := c.Param("camId")

	if err := validator.UUID(camId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	var body dto.UpdateCampaign

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign := conv.UpdateCampaign(body)

	campaign.Id = camId
	campaign.AdvertiserId = id

	needGen := false
	prompt := ""

	if body.GenText != nil {
		needGen = *body.GenText

		if body.Prompt == nil {
			resp.AbortErrMsg(c, e.New("If you use text generating, prompt is required", e.BadInput))
			return
		}
		prompt = *body.Prompt
	}

	err := ca.usecase.Update(ctx, campaign, needGen, prompt)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoCampaign(campaign)

	c.JSON(httper.StatusOK, result)
}

func (ca *Campaign) GetImage(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	camId := c.Param("camId")

	if err := validator.UUID(camId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	url, err := ca.usecase.GetImage(ctx, campaign)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := &dto.ImageUrl{
		Url: url,
	}

	c.JSON(httper.StatusOK, result)
}

func (ca *Campaign) UploadImage(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	camId := c.Param("camId")

	if err := validator.UUID(camId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	reader, err := file.Open()
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}
	defer reader.Close()

	buffer, err := io.ReadAll(reader)
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	image := &entity.Image{
		Name:        file.Filename,
		Buffer:      buffer,
		Size:        file.Size,
		ContentType: file.Header["Content-Type"][0],
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	url, uploadErr := ca.usecase.UploadImage(ctx, campaign, image)
	if uploadErr != nil {
		resp.AbortErrMsg(c, uploadErr)
		return
	}

	result := &dto.ImageUrl{
		Url: url,
	}

	c.JSON(httper.StatusOK, result)
}

func (ca *Campaign) DeleteImage(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	camId := c.Param("camId")

	if err := validator.UUID(camId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	err := ca.usecase.DeleteImage(ctx, campaign)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}

func (ca *Campaign) Delete(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	camId := c.Param("camId")

	if err := validator.UUID(camId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	campaign := &entity.Campaign{
		Id:           camId,
		AdvertiserId: id,
	}

	err := ca.usecase.Delete(ctx, campaign)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}
