package converter

import (
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/http/v1/dto"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
)

func DtoClient(client *entity.Client) *dto.Client {
	return &dto.Client{
		Id:       client.Id,
		Login:    client.Login,
		Age:      client.Age,
		Location: client.Location,
		Gender:   client.Gender,
	}
}

func BulkClient(body []dto.Client) []*entity.Client {
	result := make([]*entity.Client, 0)

	for _, client := range body {
		toAdd := &entity.Client{
			Id:       client.Id,
			Login:    client.Login,
			Age:      client.Age,
			Location: client.Location,
			Gender:   client.Gender,
		}

		result = append(result, toAdd)
	}

	return result
}
