package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("INSERT INTO users (user_name, password, access_level) VALUES ('admin', '$2a$14$/cl5uuLjw30otkYzb519QeVwUVlvbvciYQkifpyUxSddaVKDckHMW', 0);")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210430061333_create_admin_user", up, down, opts)
}
