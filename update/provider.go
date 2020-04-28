package update

// Provider is an interface that abstracts away the provider of new binaries.
// It gives you the new version number when it's available and expects it when
// you request to do the update.
// There could be new implementations with e.g. external network service providing
// new binaries.
type Provider interface {
	// IsUpdateAvailable checks if new update is available and returns its
	// version number and an error.
	IsUpdateAvailable(version int) (newversion int, err error)

	// Update performs the update to the defined version.
	Update(version int) error
}
