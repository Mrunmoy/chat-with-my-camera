package main


import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// CameraInfo is what you send to the client.
type CameraInfo struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Thumbnail string `json:"thumbnail"`
}
type App struct {
	DB     *sql.DB
	Config *Config // your config struct type
}


// NewApp sets up your App struct with DB + Config.
func NewApp(db *sql.DB, cfg *Config) *App {
	return &App{
		DB:     db,
		Config: cfg,
	}
}