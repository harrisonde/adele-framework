package postgresdriver

import (
	"database/sql"
	"fmt"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

// Build a data source name string
func BuildDSN(host, port, user, password, dbname, sslmode string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
		host, port, user, password, dbname, sslmode)
}

// Create a new Postgres builder session
func Session(pool *sql.DB) (db.Session, error) {
	session, err := postgresql.New(pool)
	if err != nil {
		return nil, err
	}

	return session, nil
}
