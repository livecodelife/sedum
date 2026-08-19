package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	// sedum:anchor:imports
)

// pingAttempts and pingDelay bound how long a service waits for its database
// at boot. A container that starts beside postgres will lose the race, so the
// standard retries rather than fataling on the first refusal - but it still
// fatals in the end, because a service that cannot reach its database has
// nothing to serve.
const (
	pingAttempts = 30
	pingDelay    = time.Second
)

func waitForDatabase(conn *sql.DB) {
	for attempt := 1; attempt <= pingAttempts; attempt++ {
		err := conn.Ping()
		if err == nil {
			return
		}
		log.Printf("database not ready (attempt %d/%d): %v", attempt, pingAttempts, err)
		time.Sleep(pingDelay)
	}
	log.Fatalf("database unreachable after %d attempts", pingAttempts)
}

func main() {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer conn.Close()

	waitForDatabase(conn)
	db.SetConn(conn)

	r := chi.NewRouter()

	// sedum:anchor:routes

	// The listen address comes from the environment, the same way the database
	// URL two dozen lines up does. A service that configures one and hardcodes
	// the other has no convention, only an accident.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, r))
}
