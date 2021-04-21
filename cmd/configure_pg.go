package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	fmt.Println(string(query), err)
	return nil
}

type nullLogger struct{}

func (d nullLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d nullLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	return nil
}

func connectToPG(prefix string) (*pg.DB, error) {
	user := config.GetString(fmt.Sprintf("%s.user", prefix))
	pass := config.GetString(fmt.Sprintf("%s.pass", prefix))
	host := config.GetString(fmt.Sprintf("%s.host", prefix))
	database := config.GetString(fmt.Sprintf("%s.database", prefix))
	port := config.GetInt(fmt.Sprintf("%s.port", prefix))
	poolSize := config.GetInt(fmt.Sprintf("%s.poolSize", prefix))
	maxRetries := config.GetInt(fmt.Sprintf("%s.maxRetries", prefix))
	timeout := config.GetInt(fmt.Sprintf("%s.connectionTimeout", prefix))

	options := &pg.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		User:       user,
		Password:   pass,
		Database:   database,
		PoolSize:   poolSize,
		MaxRetries: maxRetries,
	}
	db := pg.Connect(options)
	db.AddQueryHook(nullLogger{})
	err := waitForConnection(db, timeout)
	return db, err
}

func waitForConnection(db *pg.DB, timeout int) error {
	t := time.Duration(timeout) * time.Second
	timeoutTimer := time.NewTimer(t)
	defer timeoutTimer.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf("timed out waiting for PostgreSQL to connect")
		case <-ticker.C:
			if isConnected(db) {
				return nil
			}
		}
	}
}

func isConnected(db *pg.DB) bool {
	res, err := db.Exec("select 1")
	if err != nil {
		return false
	}
	return res.RowsReturned() == 1
}
