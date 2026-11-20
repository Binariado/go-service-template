// Package migrations embeds the SQL migration files.
//
// Naming: 000NNN_snake_case_description.up.sql and a matching .down.sql.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
