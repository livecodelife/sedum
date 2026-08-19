package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
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

// The resource's own columns, in the order its fields were declared.
//
// This list and the two below are the same invocation's three consequences: one
// addModelField fills a field, a column name, an insert argument and a scan
// target. Nothing restates a name another action already bound, so no answer can
// be consistent with itself and still disagree about spelling or order.
var {{name|receiver}}Columns = []string{
	// sedum:anchor:columns
}

// The helpers are named for the resource rather than shared, because every
// resource in a service compiles into one db package and a shared helper would
// be redeclared once per file.

func {{name|receiver}}ColumnList() string { return strings.Join({{name|receiver}}Columns, ", ") }

// {{name|receiver}}Placeholders numbers one placeholder per column, starting at
// offset - $1 for an insert, $2 for an update whose $1 is the row's id.
func {{name|receiver}}Placeholders(offset int) string {
	out := make([]string, len({{name|receiver}}Columns))
	for i := range {{name|receiver}}Columns {
		out[i] = fmt.Sprintf("$%d", i+offset)
	}
	return strings.Join(out, ", ")
}

// {{name|receiver}}Assignments writes the SET clause of an UPDATE.
func {{name|receiver}}Assignments(offset int) string {
	out := make([]string, len({{name|receiver}}Columns))
	for i, c := range {{name|receiver}}Columns {
		out[i] = fmt.Sprintf("%s = $%d", c, i+offset)
	}
	return strings.Join(out, ", ")
}

// insertArgs reads the resource's own columns off the value being written, in
// the order {{name|receiver}}Columns names them.
func (r {{name|exported}}) insertArgs() []any {
	return []any{
		// sedum:anchor:insert_args
	}
}

// scanTargets points at the fields a SELECT writes into, in the order the
// select templates list them: id, then the resource's own columns, then the
// timestamps.
func (r *{{name|exported}}) scanTargets() []any {
	return []any{
		&r.ID,
		// sedum:anchor:scan_targets
		&r.CreatedAt,
		&r.UpdatedAt,
	}
}

// sedum:anchor:queries
