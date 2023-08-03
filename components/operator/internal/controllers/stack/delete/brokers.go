package delete

import (
	"fmt"
	"reflect"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/storage/nats"
	"github.com/go-logr/logr"
)

var (
	natsClientId = "12"
)

func DeleteBrokersData(c *v1beta3.Configuration, stackName string, subjectService []string, logger logr.Logger) error {
	values := reflect.ValueOf(c.Spec.Broker)
	for i := 0; i < values.NumField(); i++ {
		switch values.Type().Field(i).Name {
		case "Nats":
			natsConfig := values.Field(i).Interface().(*v1beta3.NatsConfig)
			if natsConfig == nil {
				continue
			}

			client, err := nats.NewClient(natsConfig, natsClientId)
			defer client.Close()
			if err != nil {
				logger.Error(err, "NATS Client")
				continue
			}

			for _, service := range subjectService {
				err = nats.DeleteSubject(client, fmt.Sprintf("%s-%s", stackName, service))
				if err != nil {
					logger.Error(err, "NATS Delete subject")
					continue
				}
			}
		}

	}
	return nil
}
