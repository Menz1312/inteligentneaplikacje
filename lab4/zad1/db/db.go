package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type bdb struct {
	db *sql.DB
}

var sdb *bdb = nil

func setup(address string) error {
	if sdb != nil {
		return errors.New("setup already done")
	}
	db, err := sql.Open("mysql", address)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	sdb = &bdb{
		db: db,
	}
	return nil
}

func Query(query string, args ...any) (*sql.Rows, error) {
	return sdb.db.Query(query, args...)
}

func QueryRow(query string, args ...any) *sql.Row {
	return sdb.db.QueryRow(query, args...)
}

func Exec(query string, args ...any) (sql.Result, error) {
	return sdb.db.Exec(query, args...)
}

func Use(database string) {
	Exec(fmt.Sprintf("USE %s;", database))
}

func Prepare(query string) (*sql.Stmt, error) {
	return sdb.db.Prepare(query)
}
