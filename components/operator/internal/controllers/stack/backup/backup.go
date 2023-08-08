package backup

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/storage/pg"
	"github.com/formancehq/operator/internal/controllers/stack/storage/s3"
	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	extension = ".gz"
)

var (
	ErrCast = errors.New("cannot cast interface to string")
)

func BackupServices(c *v1beta3.Configuration, stackName string, storage s3.Storage, t *v1.Time, logger logr.Logger) error {
	values := reflect.ValueOf(c.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			if servicesValues.Type().Field(j).Name != "Postgres" {
				continue
			}

			postgresConfig, ok := servicesValues.Field(j).Interface().(v1beta3.PostgresConfig)
			if !ok {
				logger.Error(ErrCast, "cannot cast to postgresconfig")
				return ErrCast
			}

			serviceName := strings.ToLower(values.Type().Field(i).Name)
			databaseName := fmt.Sprintf("%s-%s", stackName, serviceName)

			if err := backupService(c, postgresConfig, databaseName, storage, t, logger); err != nil {
				logger.Error(err, "service not backuped"+err.Error())
				return err
			}
		}
	}
	return nil
}

func backupService(
	c *v1beta3.Configuration,
	pgConfig v1beta3.PostgresConfig,
	databaseName string,
	storage s3.Storage,
	t *v1.Time,
	logger logr.Logger,
) error {

	date := t.Format(time.RFC3339)
	fileName := fmt.Sprintf("%s-%s%s", databaseName, date, extension)

	logger.Info("backuping process " + databaseName)
	exist, err := storage.Exist(fileName)
	if err != nil {
		return err
	}

	if exist {
		logger.Info("database already backuped")
		return nil
	}

	data, err := backupPostgres(databaseName, pgConfig, logger)
	if err != nil {
		return err
	}

	err = storage.PutFile(fileName, data)
	if err != nil {
		logger.Error(err, "uploding to s3")
		return err
	}

	logger.Info("database backuped")

	return nil
}

func backupPostgres(databaseName string, conf v1beta3.PostgresConfig, logger logr.Logger) ([]byte, error) {

	data, err := pg.BackupDatabase(
		databaseName,
		conf,
	)

	if err != nil {
		logger.Error(err, "Backup database")
		return nil, err
	}

	return data, nil
}
