package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
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

var (
	ErrNotExisting = errors.New("not existing")
)

func DropDB(db *sql.DB, stackName string, serviceName string, ctx context.Context) error {
	_, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+fmt.Sprintf("\"%s-%s\"", stackName, serviceName))
	if err != nil {
		return err
	}

	return nil
}

func OpenClient(config v1beta3.PostgresConfig) (*sql.DB, error) {
	return OpenSQLDB(ConnectionOptions{
		DatabaseSourceName: config.DSN(),
		Debug:              config.Debug,
		Trace:              config.Debug,
		Writer:             os.Stdout,
		MaxIdleConns:       20,
		ConnMaxIdleTime:    time.Minute,
		MaxOpenConns:       20,
	})
}
