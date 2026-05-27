package db

import "embed"

// Migrations holds all .sql files under db/migrations.

var Seeds embed.FS //go:embed seeds/*.sql//// Seeds holds all .sql files under db/seeds.var Migrations embed.FS//go:embed migrations/*.sql//
