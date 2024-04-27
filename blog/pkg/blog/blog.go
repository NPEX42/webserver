package blog

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Blog struct {
	db *sql.DB
}

func Connect(uri string) (*Blog, error) {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		log.Printf("Failed To Open DB: %s\n", err)
		return nil, err
	}

	return &Blog{db: db}, nil
}

func (b *Blog) Close() {
	b.db.Close()
}

func (b *Blog) Query(query string, args ...any) (*sql.Rows, error) {
	return b.db.Query(query, args...)
}

func (b *Blog) Exec(query string, args ...any) (sql.Result, error) {
	tx, err := b.db.Begin()
	if err != nil {
		log.Printf("Failed To Begin Transaction: %s\n", err)
		return &EmptyResult{}, err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed To Execute Query: %q\n\t%s\n", query, err)
		return &EmptyResult{}, err
	}

	tx.Commit()

	return result, nil
}

type EmptyResult struct {
}

func (er *EmptyResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("Empty Result, Check Exec's Error!\n")
}

func (er *EmptyResult) RowsAffected() (int64, error) {
	return 0, fmt.Errorf("Empty Result, Check Exec's Error!\n")
}
