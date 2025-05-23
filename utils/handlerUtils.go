package utils

import (
	"Basic_CLI_Application/consts"
	"strings"
)

func IsStatusValid(status string) bool {
	if status == "" {
		return false
	}

	lowerStatus := strings.ToLower(status)
	return lowerStatus == consts.TodoStatusNotStarted || lowerStatus == consts.TodoStatusStarted || lowerStatus == consts.TodoStatusCompleted
}
