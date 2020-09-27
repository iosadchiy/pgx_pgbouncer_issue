package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const connstring = "host=pgbouncer  port=6432 user=postgres password=secret dbname=postgres sslmode=disable statement_cache_mode=describe"

func NewDatabasePool() (*sql.DB, error) {
	db, err := sql.Open("pgx", connstring)
	if err != nil {
		return db, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func main() {
	for {
		err := run()
		fmt.Println("run() finished: ", err)
		time.Sleep(10 * time.Second)
	}
}

func run() error {
	ctx := context.Background()

	db, err := NewDatabasePool()
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(ctx, db, wg)
	}
	wg.Wait()

	return nil
}

func worker(ctx context.Context, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		handle(ctx, db)
	}
}

func handle(ctx context.Context, db *sql.DB) {
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Millisecond)
	defer cancel()

	q := `select pg_sleep(10)`
	_, _ = db.ExecContext(ctx, q)
}
