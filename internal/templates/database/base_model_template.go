package database

func BaseModelTemplate() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base provides a UUID primary key and standard timestamps for all models.
type Base struct {
	ID        uuid.UUID      ` + "`json:\"id\" gorm:\"type:char(36);primaryKey\"`" + `
	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`json:\"deleted_at\" gorm:\"index\"`" + `
}

// BeforeCreate auto-generates a UUID if one has not been set.
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
`
}
