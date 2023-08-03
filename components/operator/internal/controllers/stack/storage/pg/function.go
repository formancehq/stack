package pg

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
)

func DropDB(db *sql.DB, stackName string, serviceName string) error {
	_, err := db.Exec("DROP DATABASE " + fmt.Sprintf("\"%s-%s\" WITH (FORCE) ", stackName, serviceName))
	if err != nil {
		return err
	}

	return nil
}

func OpenClient(config v1beta3.PostgresConfig) (*sql.DB, error) {
	debug := true
	return OpenSQLDB(ConnectionOptions{
		DatabaseSourceName: config.DSN(),
		Debug:              debug,
		Trace:              debug,
		Writer:             os.Stdout,
		MaxIdleConns:       20,
		ConnMaxIdleTime:    time.Minute,
		MaxOpenConns:       20,
	})
}
