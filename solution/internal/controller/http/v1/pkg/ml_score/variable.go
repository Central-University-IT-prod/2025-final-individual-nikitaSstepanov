package mlscore

import (
	e "github.com/nikitaSstepanov/tools/error"
	resp "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/controller/response"
)

var (
	badReqErr = e.New("Incorrect data.", e.BadInput)
	okMsg     = resp.NewMessage("Ok.")
)
