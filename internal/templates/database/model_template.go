package database

func ModelTemplate(name string) (string, error) {
	return `package models

// ` + name + ` represents the ` + name + ` domain entity.
type ` + name + ` struct {
	Base
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\" gorm:\"uniqueIndex\"`" + `
}
`, nil
}

func ModelSkeletonTemplate(name string) (string, error) {
	return `package models

// ` + name + ` is a skeleton model generated without --all.
// Add fields based on your domain needs.
type ` + name + ` struct {
	Base
}
`, nil
}
