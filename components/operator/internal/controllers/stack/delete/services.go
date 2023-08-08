package delete

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/storage/es"
	"github.com/formancehq/operator/internal/controllers/stack/storage/pg"
	"github.com/go-logr/logr"
)

func DeleteByService(c *v1beta3.Configuration, stackName string, logger logr.Logger) error {
	values := reflect.ValueOf(c.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			switch servicesValues.Type().Field(j).Name {
			case "Postgres":
				postgresConfig, ok := servicesValues.Field(j).Interface().(v1beta3.PostgresConfig)
				if !ok {
					logger.Error(ErrCast, "cannot cast to postgresconfig")
					return ErrCast
				}

				serviceName := strings.ToLower(values.Type().Field(i).Name)

				client, err := pg.OpenClient(postgresConfig)
				defer client.Close()
				if err != nil {
					logger.Error(err, "PG: Cannot open pg client")
					return err
				}

				if err := pg.DropDB(client, stackName, serviceName); err != nil {
					logger.Error(err, "PG: Error during drop")
					return err
				}

				logger.Info(fmt.Sprintf("PG: database \"%s-%s\" droped", stackName, serviceName))
			case "ElasticSearchConfig":
				elasticSearchConfig := servicesValues.Field(j).Interface().(v1beta3.ElasticSearchConfig)
				if err := es.DropESIndex(&elasticSearchConfig, logger, stackName); err != nil {
					logger.Error(err, "ELK: Error during index drop es")
					return err
				}
			}

		}
	}

	return nil
}
