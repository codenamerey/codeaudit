package utils

import "strings"

func IsConfigFlag(value string) bool {
	configTuple := strings.Split(value, "=")
	if strings.Contains(value, "--") && len(configTuple) == 2 {
		return true
	} else {
		return false
	}
}
