package mlscore

import e "github.com/nikitaSstepanov/tools/error"

const (
	mlTable = "ml_scores"
)

var (
	notFoundErr = e.New("This score wasn`t found.", e.NotFound)
)
