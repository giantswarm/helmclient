package helmclient

import "strings"

var (
	noHostErrorString = "no such host"
)

// IsNoHostError asserts noHostError.
func IsNoHostError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), noHostErrorString)
}
