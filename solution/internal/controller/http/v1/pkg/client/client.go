package client

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/httper"
	conv "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/converter"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/validator"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
	ct "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/utils/controller"
)

type Client struct {
	usecase ClientUseCase
}

func New(uc ClientUseCase) *Client {
	return &Client{
		usecase: uc,
	}
}

func (cl *Client) GetById(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	client, err := cl.usecase.GetById(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoClient(client)

	c.JSON(httper.StatusOK, result)
}

func (cl *Client) Bulk(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body []dto.Client

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	for _, client := range body {
		if err := validator.Struct(client); err != nil {
			resp.AbortErrMsg(c, err)
			return
		}
	}

	clients := conv.BulkClient(body)

	if err := cl.usecase.Bulk(ctx, clients); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusCreated, body)
}
