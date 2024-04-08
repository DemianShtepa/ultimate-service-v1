package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
)

type Config struct {
	Host     string
	User     string
	Password string
	Database string
}

func Open(config Config) (*sqlx.DB, error) {
	query := make(url.Values)
	query.Set("timezone", "utc")
	query.Set("sslmode", "disable")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.User, config.Password),
		Host:     config.Host,
		Path:     config.Database,
		RawQuery: query.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	var tmp bool

	return db.QueryRowxContext(ctx, "SELECT true").Scan(&tmp)
}
