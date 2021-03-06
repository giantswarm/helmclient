package helmclient

import (
	"errors"
	"net"
	"regexp"
	"strings"

	"github.com/giantswarm/microerror"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	resourceAlreadyExistsErrorPrefix = "rendered manifests contain a resource that already exists"
)

var resourceAlreadyExistsError = &microerror.Error{
	Kind: "resourceAlreadyExistsError",
}

// IsResourceAlreadyExists asserts resourceAlreadyExistError.
func IsResourceAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), resourceAlreadyExistsErrorPrefix) {
		return true
	}
	if c == resourceAlreadyExistsError {
		return true
	}

	return false
}

const (
	cannotReuseReleaseErrorPrefix = "cannot re-use"
)

var cannotReuseReleaseError = &microerror.Error{
	Kind: "cannotReuseReleaseError",
}

// IsCannotReuseRelease asserts cannotReuseReleaseError.
func IsCannotReuseRelease(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), cannotReuseReleaseErrorPrefix) {
		return true
	}
	if c == cannotReuseReleaseError {
		return true
	}

	return false
}

var (
	emptyChartTemplatesRegexp = regexp.MustCompile(`release \S+ failed: no objects visited`)
)

var emptyChartTemplatesError = &microerror.Error{
	Kind: "emptyChartTemplatesError",
}

// IsEmptyChartTemplates asserts emptyChartTemplatesError.
func IsEmptyChartTemplates(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if c == emptyChartTemplatesError {
		return true
	}
	if emptyChartTemplatesRegexp.MatchString(c.Error()) {
		return true
	}

	return false
}

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

const (
	invalidGZipHeaderErrorPrefix = "gzip: invalid header"
)

var invalidGZipHeaderError = &microerror.Error{
	Kind: "invalidGZipHeaderError",
}

// IsInvalidGZipHeader asserts invalidGZipHeaderError.
func IsInvalidGZipHeader(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.HasPrefix(c.Error(), invalidGZipHeaderErrorPrefix) {
		return true
	}
	if c == invalidGZipHeaderError {
		return true
	}

	return false
}

const (
	invalidManifestPrefix = "unable to build kubernetes objects from release manifest"
)

var invalidManifestError = &microerror.Error{
	Kind: "invalidManifestError",
}

// IsInvalidManifest asserts invalidManifestError.
func IsInvalidManifest(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.HasPrefix(c.Error(), invalidManifestPrefix) {
		return true
	}
	if c == invalidManifestError {
		return true
	}

	return false
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var parsingDestFailedError = &microerror.Error{
	Kind: "parsingDestFailedError",
}

// IsParsingDestFailedError asserts parsingDestFailedError.
func IsParsingDestFailedError(err error) bool {
	return microerror.Cause(err) == parsingDestFailedError
}

var parsingSrcFailedError = &microerror.Error{
	Kind: "parsingSrcFailedError",
}

// IsparsingSrcFailedError asserts parsingSrcFailedError.
func IsParsingSrcFailedError(err error) bool {
	return microerror.Cause(err) == parsingSrcFailedError
}

var pullChartFailedError = &microerror.Error{
	Kind: "pullChartFailedError",
}

// IsPullChartFailedError asserts pullChartFailedError.
func IsPullChartFailedError(err error) bool {
	return microerror.Cause(err) == pullChartFailedError
}

var pullChartNotFoundError = &microerror.Error{
	Kind: "pullChartNotFoundError",
}

// IsPullChartNotFound asserts pullChartNotFoundError.
func IsPullChartNotFound(err error) bool {
	return microerror.Cause(err) == pullChartNotFoundError
}

var pullChartTimeoutError = &microerror.Error{
	Kind: "pullChartTimeoutError",
}

// IsPullChartTimeout asserts pullChartTimeoutError.
func IsPullChartTimeout(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if c == pullChartTimeoutError {
		return true
	}

	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}

	return netErr.Timeout()
}

var releaseAlreadyExistsError = &microerror.Error{
	Kind: "releaseAlreadyExistsError",
}

// IsReleaseAlreadyExists asserts releaseAlreadyExistsError.
func IsReleaseAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if errors.Is(err, driver.ErrReleaseExists) {
		return true
	}
	if c == releaseAlreadyExistsError {
		return true
	}

	return false
}

const (
	releaseNameInvalidErrorPrefix = "invalid release name"
	releaseNameInvalidErrorSuffix = "and the length must not be longer than 53"
)

var releaseNameInvalidError = &microerror.Error{
	Kind: "releaseNameInvalidError",
}

// IsReleaseNameInvalid asserts releaseNameInvalidError.
func IsReleaseNameInvalid(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.HasPrefix(c.Error(), releaseNameInvalidErrorPrefix) {
		return true
	}
	if strings.HasSuffix(c.Error(), releaseNameInvalidErrorSuffix) {
		return true
	}
	if c == releaseNameInvalidError {
		return true
	}

	return false
}

const (
	releaseNotDeployedErrorSuffix = "has no deployed releases"
)

var releaseNotDeployedError = &microerror.Error{
	Kind: "releaseNotDeployedError",
}

// IsReleaseNotDeployed asserts releaseNotDeployedError.
func IsReleaseNotDeployed(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.HasSuffix(c.Error(), releaseNotDeployedErrorSuffix) {
		return true
	}
	if c == releaseNotDeployedError {
		return true
	}

	return false
}

var releaseNotFoundError = &microerror.Error{
	Kind: "releaseNotFoundError",
}

// IsReleaseNotFound asserts releaseNotFoundError.
func IsReleaseNotFound(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if errors.Is(err, driver.ErrReleaseNotFound) {
		return true
	}
	if c == releaseNotFoundError {
		return true
	}

	return false
}

var (
	tarballNotFoundRegexp = regexp.MustCompile(`stat \S+: no such file or directory`)
)

var tarballNotFoundError = &microerror.Error{
	Kind: "tarballNotFoundError",
}

// IsTarballNotFound asserts tarballNotFoundError.
func IsTarballNotFound(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if c == tarballNotFoundError {
		return true
	}
	if tarballNotFoundRegexp.MatchString(c.Error()) {
		return true
	}

	return false
}

var testReleaseFailureError = &microerror.Error{
	Kind: "testReleaseFailureError",
}

// IsTestReleaseFailure asserts testReleaseFailureError.
func IsTestReleaseFailure(err error) bool {
	return microerror.Cause(err) == testReleaseFailureError
}

var testReleaseTimeoutError = &microerror.Error{
	Kind: "testReleaseTimeoutError",
}

// IsTestReleaseTimeout asserts testReleaseTimeoutError.
func IsTestReleaseTimeout(err error) bool {
	return microerror.Cause(err) == testReleaseTimeoutError
}

var tooManyResultsError = &microerror.Error{
	Kind: "tooManyResultsError",
}

// IsTooManyResults asserts tooManyResultsError.
func IsTooManyResults(err error) bool {
	return microerror.Cause(err) == tooManyResultsError
}

var (
	yamlConversionFailedErrorText = "error converting YAML to JSON:"
)

var yamlConversionFailedError = &microerror.Error{
	Kind: "yamlConversionFailedError",
}

// IsYamlConversionFailed asserts yamlConversionFailedError.
func IsYamlConversionFailed(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if c == yamlConversionFailedError {
		return true
	}
	if strings.Contains(c.Error(), yamlConversionFailedErrorText) {
		return true
	}

	return false
}

var (
	validationFailedErrorText = "error validating data"
)

var validationFailedError = &microerror.Error{
	Kind: "validationFailedError",
}

// IsValidationFailedError asserts validationFailedError.
func IsValidationFailedError(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if c == validationFailedError {
		return true
	}
	if strings.Contains(c.Error(), validationFailedErrorText) {
		return true
	}

	return false
}
