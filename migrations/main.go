package main

import (
	"flag"
	"fmt"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

var (
	user     string
	pass     string
	address  string
	database string
)

func main() {
	flag.StringVar(&user, "user", "calorie-tracker-user", "user")
	flag.StringVar(&pass, "pass", "", "pass")
	flag.StringVar(&address, "address", "localhost:9000", "address")
	flag.StringVar(&database, "database", "calorie-tracker", "database")
	flag.Parse()

	db := pg.Connect(&pg.Options{
		User:     user,
		Password: pass,
		Database: database,
		Addr:     address,
	})

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		panic(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}
