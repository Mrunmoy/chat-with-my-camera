package main


import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB     *sql.DB
	Config Config // your config struct type
}
