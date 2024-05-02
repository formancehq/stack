package modules

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/formancehq/payments/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
)

var Payments = internal.NewModule("payments").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("paymentsapi", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"api",
					"serve",
					"--auth-enabled=false",
					"--postgres-uri=" + test.GetDatabaseSourceName("payments"),
					"--config-encryption-key=encryption-key",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=payments",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--publisher-topic-mapping=*:%s-payments", test.ID()),
					"--listen=0.0.0.0:0",
					"--auto-migrate",
				}
			}).WithRoutingFunc("payments", func(path, method string) bool {
			switch {
			case strings.HasPrefix(path, "/payments"):
				return true
			case strings.HasPrefix(path, "/accounts"):
				return method == http.MethodGet || method == http.MethodPost
			case strings.HasPrefix(path, "/bank-accounts"),
				strings.HasPrefix(path, "/transfer-initiations"):
				return method == http.MethodGet
			default:
				return false
			}
		}),
		internal.NewCommandService("paymentsconnectors", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"connectors",
					"serve",
					"--auth-enabled=false",
					"--postgres-uri=" + test.GetDatabaseSourceName("payments"),
					"--config-encryption-key=encryption-key",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=payments",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--publisher-topic-mapping=*:%s-payments", test.ID()),
					"--listen=0.0.0.0:0",
					"--auto-migrate",
				}
			}).WithRoutingFunc("payments", func(path, method string) bool {
			switch {
			case strings.HasPrefix(path, "/bank-accounts"):
				return method == http.MethodPost
			case strings.HasPrefix(path, "/transfer-initiations"):
				return method == http.MethodPost || method == http.MethodDelete
			case strings.HasPrefix(path, "/connectors"):
				return true
			default:
				return false
			}
		}),
	)
