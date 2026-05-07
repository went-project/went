package handlers

// CreateAllRelated generates all related templates in full mode.
func CreateAllRelated(name string) error {
	return generateRelated(name, true)
}

// CreateAllRelatedSkeleton generates all related templates in skeleton mode.
func CreateAllRelatedSkeleton(name string) error {
	return generateRelated(name, false)
}
