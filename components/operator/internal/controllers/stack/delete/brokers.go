package delete

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/storage/nats"
	"github.com/go-logr/logr"
)

var (
	natsClientId = "12"
)
var (
	ErrCast = errors.New("cannot cast interface to string")
)

func DeleteByBrokers(c *v1beta3.Configuration, stackName string, subjectService []string, logger logr.Logger) error {
	values := reflect.ValueOf(c.Spec.Broker)
	for i := 0; i < values.NumField(); i++ {
		switch values.Type().Field(i).Name {
		case "Nats":
			natsConfig := values.Field(i).Interface().(*v1beta3.NatsConfig)
			if natsConfig == nil {
				continue
			}
			if err := deleleNatsSubjects(natsConfig, stackName, subjectService, logger); err != nil {
				return err
			}
		}

	}
	return nil
}

func deleleNatsSubjects(config *v1beta3.NatsConfig, stackName string, subjectService []string, logger logr.Logger) error {
	client, err := nats.NewClient(config, natsClientId)
	if err != nil {
		logger.Error(err, "NATS: client")
		return err
	}
	defer client.Close()

	for _, service := range subjectService {
		stackSubjectName := fmt.Sprintf("%s-%s", stackName, service)

		exist, err := nats.ExistSubject(client, stackSubjectName)
		if err != nil {
			logger.Error(err, "NATS: subject existancy check")
			return err
		}

		// It meens it has already been deleted, and just not exists anymore
		if !exist {
			continue
		}

		// Delete subject when it exists
		err = nats.DeleteSubject(client, stackSubjectName)
		if err != nil {
			logger.Error(err, "NATS: delete subject")
			return err
		}
	}

	return nil
}
