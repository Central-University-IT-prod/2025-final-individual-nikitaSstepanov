package telegram

import "fmt"

func stateQuery(id uint64) string {
	return fmt.Sprintf("state:%d", id)
}

func sessionQuery(id uint64) string {
	return fmt.Sprintf("session:%d", id)
}

func newQuery(id uint64) string {
	return fmt.Sprintf("new:%d", id)
}

func campaignQuery(id uint64) string {
	return fmt.Sprintf("campaign:%d", id)
}
