package handlers

// CreateControllerHelper ensures internal/helpers/helpers.go exists for the project.
// Delegates to CreateHelpers for the actual file creation.
func CreateControllerHelper(basePath string) error {
	return CreateHelpers(basePath)
}
