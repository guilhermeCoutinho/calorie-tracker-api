package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`

		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_name VARCHAR NOT NULL,
			password VARCHAR NOT NULL,
			access_level VARCHAR,
			calorie_limit INTEGER NOT NULL DEFAULT 0,

			created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
		  );
		  
		  CREATE UNIQUE INDEX IF NOT EXISTS user_name ON users (user_name);
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE users;")
		return err
	}

	opts := migrations.MigrationOptions{}
	migrations.Register("20210420020046_create_users_table", up, down, opts)
}
