// +build k8srequired

package upgradetiller

import (
	"github.com/giantswarm/microerror"
)

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var tooManyResultsError = &microerror.Error{
	Kind: "tooManyResultsError",
}

// IsTooManyResults asserts tooManyResultsError.
func IsTooManyResults(err error) bool {
	return microerror.Cause(err) == tooManyResultsError
}
