package database

func ModelTemplate(name string) (string, error) {
	return `package models

import "gorm.io/gorm"

type ` + name + ` struct {
	gorm.Model
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\" gorm:\"uniqueIndex\"`" + `
}
`, nil
}

func ModelSkeletonTemplate(name string) (string, error) {
	return `package models

import "gorm.io/gorm"

// ` + name + ` is a skeleton model generated without --all.
// Add fields based on your domain needs.
type ` + name + ` struct {
	gorm.Model
}
`, nil
}
