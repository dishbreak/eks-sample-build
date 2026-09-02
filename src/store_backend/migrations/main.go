package migrations

import "embed"

//go:embed sql/*.sql
var EmbeddedMigrations embed.FS
