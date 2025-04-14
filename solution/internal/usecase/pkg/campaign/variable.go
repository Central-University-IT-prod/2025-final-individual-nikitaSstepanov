package campaign

import e "github.com/nikitaSstepanov/tools/error"

var (
	forbiddenErr = e.New("Forbidden.", e.Forbidden)
	badDateErr   = e.New("Invalid start or end dates", e.BadInput)
	badReqErr    = e.New("Invalid data. Campaign has already started", e.BadInput)
	badAgeErr    = e.New("Invalid data.", e.BadInput)
)
