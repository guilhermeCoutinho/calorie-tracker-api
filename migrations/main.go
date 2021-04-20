package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

const directory = "./"

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

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
