package db

import (
	"context"
	"database/sql"
	"time"
	// sedum:anchor:imports
)

// conn is the pool this package queries through. main wires it at startup.
var conn *sql.DB

// SetConn wires the pool.
func SetConn(db *sql.DB) { conn = db }

// ErrNotFound is what a query returns when it matches no row. Handlers map it
// to 404 without importing database/sql themselves.
var ErrNotFound = sql.ErrNoRows

// Ready reports whether the pool is reachable. It also keeps the context
// import used by this file itself rather than only by injected queries.
func Ready(ctx context.Context) error { return conn.PingContext(ctx) }

// One resource per file, named for the file. Every table in this standard
// carries the same three columns; everything between them is the service's.
type {{name|exported}} struct {
	ID int64 `json:"id"`

	// sedum:anchor:fields

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// sedum:anchor:queries
