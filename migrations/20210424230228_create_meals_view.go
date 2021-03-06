package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
		CREATE VIEW meals_with_limit (
			id,
			user_id,
			meal,
			calories,
			date,
			time_seconds,
			total_calories_for_day,
			above_limit,
			created_at,
			updated_at
		)
		AS SELECT
			id,
			user_id,
			meal,
			calories,
			date,
			time_seconds,
			(SELECT sum(calories) AS sum FROM meals meals_1 WHERE meals_1.user_id=meals.user_id AND meals_1.date=meals.date),
			(SELECT (CASE WHEN sum > calorie_limit THEN TRUE ELSE FALSE END) AS above_limit FROM (select p.sum, q.calorie_limit FROM (SELECT sum(calories) AS sum FROM meals meals_1 WHERE meals_1.user_id=meals.user_id AND meals_1.date=meals.date) AS p, (select calorie_limit as calorie_limit from users WHERE id=user_id) AS q) AS r),
			created_at,
			updated_at
		from meals;
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210424230228_create_meals_view", up, down, opts)
}
