package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDatabasePool(ctx context.Context) (*pgxpool.Pool, error) {
	s := "host=pgbouncer port=6432 user=postgres password=secret dbname=postgres pool_max_conns=10"
	c, err := pgxpool.ParseConfig(s)
	if err != nil {
		panic("failed to parse postgres config: " + err.Error())
	}

	c.MaxConns = 10
	c.ConnConfig.TLSConfig = nil

	return pgxpool.ConnectConfig(ctx, c)
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

	db, err := NewDatabasePool(ctx)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(ctx, db, wg)
	}
	wg.Wait()

	return nil
}

func worker(ctx context.Context, db *pgxpool.Pool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		handle(ctx, db)
	}
}

func handle(ctx context.Context, db *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Millisecond)
	defer cancel()

	q := `select pg_sleep(10)`
	rows, _ := db.Query(ctx, q)
	defer rows.Close()
}
