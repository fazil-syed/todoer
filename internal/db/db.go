package db

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"time"
)

var DB *sql.DB

func Init(ctx context.Context) {
	db, err := sql.Open("sqlite", "todoer.db")
	if err != nil {
		log.Fatal(err)
	}
	timeOutCtx, timeOutcancel := context.WithTimeout(ctx, 5*time.Second)
	defer timeOutcancel()
	if err := db.PingContext(timeOutCtx); err != nil {
		log.Fatal(err)
	}
	DB = db

}

//go:embed schema.sql
var schema string

func Migrate(ctx context.Context) error {
	_, err := DB.ExecContext(ctx, schema)
	if err != nil {
		return err
	}
	return nil
}
