package handlers

// CreateAllRelated generates all related templates in full mode.
func CreateAllRelated(name string) error {
	if err := CreateModel(name, true); err != nil {
		return err
	}
	if err := CreateMigration(name, true); err != nil {
		return err
	}
	if err := CreateController(name, true); err != nil {
		return err
	}
	if err := CreateRouter(name, true); err != nil {
		return err
	}
	if err := CreateResource(name, true); err != nil {
		return err
	}
	return nil
}

// CreateAllRelatedSkeleton generates all related templates in skeleton mode.
func CreateAllRelatedSkeleton(name string) error {
	if err := CreateModel(name, false); err != nil {
		return err
	}
	if err := CreateMigration(name, false); err != nil {
		return err
	}
	if err := CreateController(name, false); err != nil {
		return err
	}
	if err := CreateRouter(name, false); err != nil {
		return err
	}
	if err := CreateResource(name, false); err != nil {
		return err
	}
	return nil
}
