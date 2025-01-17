package models

import (
	"strings"
)

func fromNameToID(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
