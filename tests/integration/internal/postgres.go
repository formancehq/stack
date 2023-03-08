package internal

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/gomega"
	"github.com/xo/dburl"
)

func getPostgresDSNString() string {
	if fromEnv := os.Getenv("POSTGRES_DSN"); fromEnv != "" {
		return fromEnv
	}
	return "postgres://formance:formance@localhost:5432/formance?sslmode=disable"
}

func getPostgresDSN() (*dburl.URL, error) {
	return dburl.Parse(getPostgresDSNString())
}

func createDatabases() {
	conn, err := pgx.Connect(ctx, getPostgresDSNString())
	Expect(err).To(BeNil())
	for _, component := range []string{"ledger", "wallets", "orchestration", "auth", "payments", "webhooks"} {
		_, err = conn.Exec(ctx, databaseNameForComponent(component))
		Expect(err).To(BeNil())
	}
}

func databaseNameForComponent(name string) string {
	return fmt.Sprintf(`CREATE DATABASE "%s-%s";`, actualTestID, name)
}
