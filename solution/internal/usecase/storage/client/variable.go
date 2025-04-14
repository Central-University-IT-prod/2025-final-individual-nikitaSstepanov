package client

import (
	e "github.com/nikitaSstepanov/tools/error"
)

const (
	clientsTable = "clients"
)

var (
	notFoundErr = e.New("This client wasn`t found.", e.NotFound)
)
