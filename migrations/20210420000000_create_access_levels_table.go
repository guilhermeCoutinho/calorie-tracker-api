package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS access_levels (
			level int PRIMARY KEY,
			name VARCHAR
		);

		INSERT INTO access_levels (level, name) VALUES (0, 'admin');
		INSERT INTO access_levels (level, name) VALUES (1, 'manager');
		INSERT INTO access_levels (level, name) VALUES (2, 'regular_user');
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210420000000_create_access_levels_table", up, down, opts)
}
