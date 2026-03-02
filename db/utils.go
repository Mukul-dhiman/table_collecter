package db

import (
	"fmt"
	"regexp"
	"strings"
)

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func SanitizeColumnName(name string, index int) string {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Sprintf("col_%d", index)
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = regexp.MustCompile(`[^a-z0-9_]`).ReplaceAllString(name, "")

	if name == "" {
		return fmt.Sprintf("col_%d", index)
	}

	return name
}
