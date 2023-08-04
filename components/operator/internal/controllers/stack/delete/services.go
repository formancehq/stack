package delete

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/storage/es"
	"github.com/formancehq/operator/internal/controllers/stack/storage/pg"
	"github.com/formancehq/operator/internal/controllers/stack/storage/s3"
	"github.com/go-logr/logr"
)

const (
	extension = ".gz"
)

func DeleteServiceData(c *v1beta3.Configuration, stackName string, logger logr.Logger) error {
	values := reflect.ValueOf(c.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			switch servicesValues.Type().Field(j).Name {
			case "Postgres":
				postgresConfig, ok := servicesValues.Field(j).Interface().(v1beta3.PostgresConfig)
				if !ok {
					continue
				}

				serviceName := strings.ToLower(values.Type().Field(i).Name)

				client, err := pg.OpenClient(postgresConfig)
				defer client.Close()
				if err != nil {
					logger.Error(err, "PG: Cannot open pg client")
					continue
				}

				if err := pg.DropDB(client, stackName, serviceName); err != nil {
					logger.Error(err, "PG: Error during drop")
					continue
				}

				logger.Info(fmt.Sprintf("PG: database \"%s-%s\" droped", stackName, serviceName))
			case "ElasticSearchConfig":
				elasticSearchConfig := servicesValues.Field(j).Interface().(v1beta3.ElasticSearchConfig)
				if err := es.DropESIndex(&elasticSearchConfig, logger, stackName); err != nil {
					logger.Error(err, "ELK: Error during index drop es")
					continue
				}
			}

		}
	}

	return nil
}

func BackupServicesData(c *v1beta3.Configuration, stack *v1beta3.Stack, storage s3.Storage, logger logr.Logger) error {
	stackName := stack.Name
	values := reflect.ValueOf(c.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			if servicesValues.Type().Field(j).Name != "Postgres" {
				continue
			}

			postgresConfig, ok := servicesValues.Field(j).Interface().(v1beta3.PostgresConfig)
			if !ok {
				logger.Error(fmt.Errorf("%s", "CAST"), "cannot cast to postgresconfig")
				continue
			}

			serviceName := strings.ToLower(values.Type().Field(i).Name)

			databaseName := fmt.Sprintf(
				"%s-%s",
				stackName,
				serviceName,
			)

			logger.Info("backuping " + databaseName)

			if err := backupPostgres(
				databaseName,
				postgresConfig,
				storage,
				logger,
			); err != nil {
				logger.Error(err, "database not backuped"+err.Error())
				continue
			}

			logger.Info("database backuped")

		}
	}
	return nil
}

func backupPostgres(databaseName string, conf v1beta3.PostgresConfig, storage s3.Storage, logger logr.Logger) error {
	date := time.Now().Format(time.RFC3339)

	data, err := pg.BackupDatabase(
		databaseName,
		conf,
	)

	if err != nil {
		logger.Error(err, "Backup database")
		return err
	}

	err = storage.PutFile(fmt.Sprintf("%s-%s%s", databaseName, date, extension), data)
	if err != nil {
		logger.Error(err, "Uploding to s3")
		return err
	}

	return nil
}
