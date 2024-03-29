package helmclient

import "time"

// Chart returns information about a Helm Chart.
type Chart struct {
	// Annotations is map of key:value pairs set by Helm Chart
	// maintainers
	Annotations map[string]string
	// Version is the version of the Helm Chart.
	Version string
}

// ReleaseContent returns status information about a Helm Release.
type ReleaseContent struct {
	// AppVersion is the app version of the Helm Chart that has been deployed.
	AppVersion string
	// Description is a human-friendly "log entry" about this Helm release.
	Description string
	// LastDeployed is the time the Helm Chart was last deployed.
	LastDeployed time.Time
	// Name is the name of the Helm Release.
	Name string
	// Revision is the revision number of the Helm Release.
	Revision int
	// Status is the Helm status code of the Release.
	Status string
	// Values are the values provided when installing the Helm Release.
	Values map[string]interface{}
	// Version is the version of the Helm Chart that has been deployed.
	Version string
}

// ReleaseHistory returns version information about a Helm Release.
type ReleaseHistory struct {
	// AppVersion is the app version of the Helm Chart that has been deployed.
	AppVersion string
	// Description is a human-friendly "log entry" about this Helm release.
	Description string
	// LastDeployed is the time the Helm Chart was last deployed.
	LastDeployed time.Time
	// Name is the name of the Helm Release.
	Name string
	// Revision is the revision number of the Helm Release.
	Revision int
	// Status is the Helm status code of the Release.
	Status string
	// Version is the version of the Helm Chart that has been deployed.
	Version string
}
