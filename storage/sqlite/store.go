package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/numary/ledger/core"
	"github.com/numary/ledger/ledger/query"
	"log"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

//go:embed migration
var migrations embed.FS

type SQLiteStore struct {
	ledger string
	db     *sql.DB
}

func NewStore(name string) (*SQLiteStore, error) {
	dbpath := fmt.Sprintf(
		"file:%s?_journal=WAL",
		path.Join(
			viper.GetString("storage.dir"),
			fmt.Sprintf(
				"%s_%s.db",
				viper.GetString("storage.sqlite.db_name"),
				name,
			),
		),
	)

	log.Printf("opening %s\n", dbpath)

	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		return nil, err
	}

	return &SQLiteStore{
		ledger: name,
		db:     db,
	}, nil
}

func (s *SQLiteStore) Name() string {
	return s.ledger
}

func (s *SQLiteStore) LastTransaction() (*core.Transaction, error) {
	var lastTransaction core.Transaction

	q := query.New()
	q.Modify(query.Limit(1))

	c, err := s.FindTransactions(q)
	if err != nil {
		return nil, err
	}

	txs := (c.Data).([]core.Transaction)
	if len(txs) > 0 {
		lastTransaction = txs[0]
		return &lastTransaction, nil
	}
	return nil, nil
}

func (s *SQLiteStore) LastMetaID() (int64, error) {
	count, err := s.CountMeta()
	if err != nil {
		return 0, err
	}
	return count - 1, nil
}

func (s *SQLiteStore) Initialize() error {
	log.Println("initializing sqlite db")

	statements := []string{}

	entries, err := migrations.ReadDir("migration")

	if err != nil {
		return err
	}

	for _, m := range entries {
		log.Printf("running migration %s\n", m.Name())

		b, err := migrations.ReadFile(path.Join("migration", m.Name()))

		if err != nil {
			return err
		}

		plain := strings.ReplaceAll(string(b), "VAR_LEDGER_NAME", s.ledger)

		statements = append(
			statements,
			strings.Split(plain, "--statement")...,
		)
	}

	for i, statement := range statements {
		_, err = s.db.Exec(
			statement,
		)

		if err != nil {
			fmt.Println(err)
			err = fmt.Errorf("failed to run statement %d: %w", i, err)
			return err
		}
	}

	return nil
}

func (s *SQLiteStore) Close() error {
	log.Println("sqlite db closed")
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}
