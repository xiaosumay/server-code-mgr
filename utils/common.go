package utils

import (
	"regexp"
	"strings"
)

func DefaultValue(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

var (
	pattern *regexp.Regexp
	Repositories = make(map[string]Repo)
)

func init() {
	pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)
}

// Quote returns a shell-escaped version of the string s. The returned value
// is a string that can safely be used as one token in a shell command line.
func Quote(s string) string {
	if len(s) == 0 {
		return ""
	}
	if pattern.MatchString(s) {
		return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
	}

	return s
}
