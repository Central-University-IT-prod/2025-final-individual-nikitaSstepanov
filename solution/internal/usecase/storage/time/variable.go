package time

import (
	e "github.com/nikitaSstepanov/tools/error"
)

const (
	redisExpires = 0
)

var (
	notFoundErr = e.New("Time wasn`t found.", e.NotFound)
)
