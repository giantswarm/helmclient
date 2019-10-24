package helmclient

import "strings"

// isNoHostError asserts no route to Host error.
func isNoHostError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "no such host")
}
