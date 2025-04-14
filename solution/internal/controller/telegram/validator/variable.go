package validator

import (
	e "github.com/nikitaSstepanov/tools/error"
)

type Arg int

var (
	lenErr  = e.New("Bad string length", e.BadInput)
	uuidErr = e.New("Id must be uuid", e.BadInput)
)

type uuid struct {
	Value string `validate:"uuid"`
}
