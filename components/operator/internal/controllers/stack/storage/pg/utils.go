package pg

import (
	"database/sql"
	"io"
	"time"

	_ "github.com/lib/pq"
)

type ConnectionOptions struct {
	DatabaseSourceName string
	Debug              bool
	Trace              bool
	Writer             io.Writer
	MaxIdleConns       int
	MaxOpenConns       int
	ConnMaxIdleTime    time.Duration
}

func OpenSQLDB(options ConnectionOptions) (*sql.DB, error) {
	sqldb, err := sql.Open("postgres", options.DatabaseSourceName)
	if err != nil {
		return nil, err
	}
	if options.MaxIdleConns != 0 {
		sqldb.SetMaxIdleConns(options.MaxIdleConns)
	}
	if options.ConnMaxIdleTime != 0 {
		sqldb.SetConnMaxIdleTime(options.ConnMaxIdleTime)
	}
	if options.MaxOpenConns != 0 {
		sqldb.SetMaxOpenConns(options.MaxOpenConns)
	}

	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	return sqldb, nil
}
