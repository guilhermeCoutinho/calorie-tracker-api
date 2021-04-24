package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS meals (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users (id),
			meal VARCHAR NOT NULL,
			calories INTEGER NOT NULL,
			DATE TIMESTAMP,
			above_limit BOOL,

			created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
		  );
		  
		  CREATE INDEX IF NOT EXISTS user_id ON meals (user_id);
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("drop table meals;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210424170501_create_meals_table", up, down, opts)
}
