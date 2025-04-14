package campaign

import e "github.com/nikitaSstepanov/tools/error"

const (
	campaignTable     = "campaigns"
	targetTable       = "targeting"
	billingTable      = "billing"
	dailyBillingTable = "daily_billing"
	impressionsTable  = "impressions"
	clicksTable       = "clicks"
)

var (
	badBillingErr = e.New("Billing is required.", e.BadInput)
	notFoundErr   = e.New("This campaign wasn`t found.", e.NotFound)
)
