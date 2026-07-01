package database

import "went/internal/utils"

func MigrationTemplates(name string) (string, string, error) {
	tableName := utils.Pluralize(name)

	up := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
	id CHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_` + tableName + `_deleted_at ON ` + tableName + ` (deleted_at);
`

	down := `DROP TABLE IF EXISTS ` + tableName + `;
`

	return up, down, nil
}

func MigrationSkeletonTemplates(name string) (string, string, error) {
	tableName := utils.Pluralize(name)

	up := `-- TODO: write your up migration for ` + tableName + `
-- Example:
-- CREATE TABLE IF NOT EXISTS ` + tableName + ` (
--   id CHAR(36) PRIMARY KEY
-- );
`

	down := `-- TODO: write your down migration for ` + tableName + `
-- Example:
-- DROP TABLE IF EXISTS ` + tableName + `;
`

	return up, down, nil
}
