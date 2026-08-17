package db

import (
	"context"
	"testing"
)

// requireConn skips rather than fails when no database is configured, so the
// suite stays runnable on a machine without one, and hands back the context
// the queries take.
func requireConn(t *testing.T) context.Context {
	t.Helper()
	if conn == nil {
		t.Skip("no database connection configured")
	}
	return context.Background()
}

// sedum:anchor:tests
