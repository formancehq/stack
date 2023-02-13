package testing

import (
	"math/rand"

	componentsv1beta3 "github.com/formancehq/operator/apis/components/v1beta3"
	v1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	apisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDumbVersions() *v1beta3.Versions {
	return &v1beta3.Versions{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: v1beta3.VersionsSpec{
			Control:        uuid.NewString(),
			Ledger:         uuid.NewString(),
			Payments:       uuid.NewString(),
			Search:         uuid.NewString(),
			Auth:           uuid.NewString(),
			Webhooks:       uuid.NewString(),
			Wallets:        uuid.NewString(),
			Counterparties: uuid.NewString(),
			Orchestration:  uuid.NewString(),
		},
	}
}

func NewDumbConfiguration() *v1beta3.Configuration {
	return &v1beta3.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: v1beta3.ConfigurationSpec{
			Services: v1beta3.ConfigurationServicesSpec{
				Auth: v1beta3.AuthSpec{
					Postgres: NewDumpPostgresConfig(),
				},
				Control: v1beta3.ControlSpec{},
				Ledger: v1beta3.LedgerSpec{
					Postgres: NewDumpPostgresConfig(),
				},
				Payments: v1beta3.PaymentsSpec{
					Postgres: NewDumpPostgresConfig(),
				},
				Search: v1beta3.SearchSpec{
					ElasticSearchConfig: NewDumpElasticSearchConfig(),
				},
				Webhooks: v1beta3.WebhooksSpec{
					Postgres: NewDumpPostgresConfig(),
				},
				Wallets: v1beta3.WalletsSpec{},
				Counterparties: v1beta3.CounterpartiesSpec{
					Postgres: NewDumpPostgresConfig(),
				},
			},
			Broker: NewDumbBrokerConfig(),
		},
	}
}

func NewDumpKafkaConfig() apisv1beta2.KafkaConfig {
	return apisv1beta2.KafkaConfig{
		Brokers: []string{"kafka:1234"},
	}
}

func NewDumbBrokerConfig() componentsv1beta3.Broker {
	return componentsv1beta3.Broker{
		Kafka: func() *apisv1beta2.KafkaConfig {
			ret := NewDumpKafkaConfig()
			return &ret
		}(),
	}
}

func NewDumpElasticSearchConfig() componentsv1beta3.ElasticSearchConfig {
	return componentsv1beta3.ElasticSearchConfig{
		Scheme: "http",
		Host:   "elasticsearch",
		Port:   9200,
	}
}

func NewDumpPostgresConfig() apisv1beta2.PostgresConfig {
	return apisv1beta2.PostgresConfig{
		Port:     5432,
		Host:     "postgres",
		Username: "username",
		Password: "password",
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func NewStackName() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
