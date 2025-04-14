package advertiser

import e "github.com/nikitaSstepanov/tools/error"

const (
	advertisersTable = "advertisers"
)

var (
	notFoundErr = e.New("This advertiser wasn`t found.", e.NotFound)
)
