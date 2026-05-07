package handlers

type GeneratorType string

const (
	GeneratorTypeModel      GeneratorType = "model"
	GeneratorTypeMigration  GeneratorType = "migration"
	GeneratorTypeController GeneratorType = "controller"
	GeneratorTypeRouter     GeneratorType = "router"
	GeneratorTypeResource   GeneratorType = "resource"
	GeneratorTypeRequest    GeneratorType = "request"
)

type GenerationOptions struct {
	Name             string
	Type             GeneratorType
	All              bool
	Force            bool
	DryRun           bool
	IncludeMigration bool
}

func Generate(opts GenerationOptions) error {
	if opts.DryRun {
		return nil
	}

	if opts.All {
		return generateRelated(opts.Name, true)
	}

	switch opts.Type {
	case GeneratorTypeModel:
		if err := CreateModel(opts.Name, false); err != nil {
			return err
		}
		if opts.IncludeMigration {
			return CreateMigration(opts.Name, false)
		}
		return nil
	case GeneratorTypeMigration:
		return CreateMigration(opts.Name, false)
	case GeneratorTypeController:
		return CreateController(opts.Name, false)
	case GeneratorTypeRouter:
		return CreateRouter(opts.Name, false)
	case GeneratorTypeResource:
		return CreateResource(opts.Name, false)
	case GeneratorTypeRequest:
		return CreateRequest(opts.Name, false)
	default:
		return nil
	}
}

func generateRelated(name string, full bool) error {
	if full {
		return generateFullRelated(name)
	}
	return generateSkeletonRelated(name)
}

func generateFullRelated(name string) error {
	if err := CreateModel(name, true); err != nil {
		return err
	}
	if err := CreateMigration(name, true); err != nil {
		return err
	}
	if err := CreateControllerHelper("."); err != nil {
		return err
	}
	if err := CreateRequest(name, true); err != nil {
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

func generateSkeletonRelated(name string) error {
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
