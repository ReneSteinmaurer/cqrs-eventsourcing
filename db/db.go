package db

import (
	"context"
	config2 "cqrs-playground/config"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	Pool   *pgxpool.Pool
	ctx    context.Context
	cancel context.CancelFunc
}

func NewDB(ctx context.Context) *DB {
	ctx, cancel := context.WithCancel(ctx)
	db := &DB{
		ctx:    ctx,
		cancel: cancel,
	}
	db.connect()
	return db
}

func (d *DB) connect() {
	fmt.Println("Verbinden mit Datenbank..")
	config, err := config2.LoadConfig()
	if err != nil {
		d.Shutdown()
		log.Fatal(err)
	}
	fmt.Printf("Connection string=%s\n", config.Database.SafeConnectionString())

	pool, err := pgxpool.New(context.Background(), config.Database.ConnectionString())
	if err != nil {
		d.Shutdown()
		log.Fatal(err)
	}

	d.Pool = pool
	err = pool.Ping(context.Background())
	if err != nil {
		d.Shutdown()
		log.Fatal(err)
	}

	err = StartMigration(config.Database.ConnectionString())
	if err != nil {
		d.Shutdown()
		log.Fatal(err)
	}

	go func() {
		select {
		case <-d.ctx.Done():
			d.Shutdown()
		}
	}()

	fmt.Println("Erfolgreich mit der Datenbank verbunden!")
}

func (d *DB) Shutdown() {
	d.cancel()
	if d.Pool != nil {
		d.Pool.Close()
		d.Pool = nil
	}
}
