package wallets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
)

type module struct{}

func (w module) Name() string {
	return "wallets"
}

func (w module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: service,
		},
		"v0.4.3": {
			Services: service,
			PostUpgrade: func(ctx context.Context, upgrader modules.JobRunner, config modules.ReconciliationConfig) (bool, error) {

				_, ok := config.Stack.Status.Ports["ledger"]
				if !ok {
					return false, errors.New("not ready, missing ledger port")
				}
				_, ok = config.Stack.Status.Ports["ledger"]["ledger"]
				if !ok {
					return false, errors.New("not ready, missing ledger port")
				}

				accounts, err := api.FetchAllPaginated[account](ctx, http.DefaultClient, ledgerUrl(config)+"/accounts", url.Values{})
				if err != nil {
					return false, errors.Wrap(err, "fetching accounts")
				}

				for _, account := range accounts {
					walletCustomMetadata, ok := account.Metadata["wallets/custom_data"]
					if ok && walletCustomMetadata != "" {
						if err := updateAccountMetadataForLedgerV2(ctx, config, account); err != nil {
							return false, errors.Wrapf(err, "updating account metadata of account: %s", account.Address)
						}
					}
				}

				return true, nil
			},
		},
		"v0.4.4": {
			Services: service,
			PostUpgrade: func(ctx context.Context, upgrader modules.JobRunner, config modules.ReconciliationConfig) (bool, error) {

				_, ok := config.Stack.Status.Ports["ledger"]
				if !ok {
					return false, errors.New("not ready, missing ledger port")
				}

				_, ok = config.Stack.Status.Ports["ledger"]["ledger"]
				if !ok {
					return false, errors.New("not ready, missing ledger port")
				}

				accounts, err := api.FetchAllPaginated[account](ctx, http.DefaultClient, ledgerUrl(config)+"/accounts", url.Values{})
				if err != nil {
					return false, errors.Wrap(err, "fetching accounts")
				}

				for _, account := range accounts {
					customData := map[string]any{}
					updated := false
					for k, v := range account.Metadata {
						switch {
						case strings.HasPrefix(k, "wallets/custom_data_"):
							customData[strings.TrimPrefix(k, "wallets/custom_data_")] = v
							delete(account.Metadata, k)
							updated = true
						case k == "destination" || k == "void_destination" || k == "wallets/holds/subject":
							switch v := v.(type) {
							case string:
								m := make(map[string]any)
								if err := json.Unmarshal([]byte(v), &m); err != nil {
									return false, err
								}
								account.Metadata[k] = m
								updated = true
							}
						}
					}
					if !updated {
						continue
					}

					account.Metadata["wallets/custom_data"] = customData

					if err := updateMetadata(ctx, config, account); err != nil {
						return false, err
					}
				}

				return true, nil
			},
		},
	}
}

var Module = &module{}

var _ modules.Module = Module

func init() {
	modules.Register(Module)
}

func service(ctx modules.ReconciliationConfig) modules.Services {
	return modules.Services{{
		HasVersionEndpoint: true,
		ExposeHTTP:         modules.DefaultExposeHTTP,
		ListenEnvVar:       "LISTEN",
		Annotations:        ctx.Configuration.Spec.Services.Wallets.Annotations.Service,
		AuthConfiguration: func(config modules.ServiceInstallConfiguration) stackv1beta3.ClientConfiguration {
			return stackv1beta3.NewClientConfiguration()
		},
		Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
			return modules.Container{
				Env: modules.ContainerEnv{
					modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
					modules.Env("STACK_CLIENT_ID", resolveContext.Stack.Status.StaticAuthClients["wallets"].ID),
					modules.Env("STACK_CLIENT_SECRET", resolveContext.Stack.Status.StaticAuthClients["wallets"].Secrets[0]),
				},
				Image: modules.GetImage("wallets", resolveContext.Versions.Spec.Wallets),
				Resources: modules.GetResourcesWithDefault(
					resolveContext.Configuration.Spec.Services.Wallets.ResourceProperties,
					modules.ResourceSizeSmall(),
				),
			}
		},
	}}
}

type account struct {
	Address  string         `json:"address"`
	Metadata map[string]any `json:"metadata"`
}

func ledgerUrl(ctx modules.ReconciliationConfig) string {
	return fmt.Sprintf("http://ledger.%s.svc.cluster.local:%d/wallets-002",
		ctx.Stack.Name,
		ctx.Stack.Status.Ports["ledger"]["ledger"])
}

func updateMetadata(ctx context.Context, config modules.ReconciliationConfig, account account) error {
	data, err := json.Marshal(account.Metadata)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ledgerUrl(config)+"/accounts/"+account.Address+"/metadata", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d while waiting for %d or %d", rsp.StatusCode, http.StatusOK, http.StatusNoContent)
	}

	return nil
}

func updateAccountMetadataForLedgerV2(ctx context.Context, config modules.ReconciliationConfig, account account) error {

	customMetadataRaw := account.Metadata["wallets/custom_data"]
	newMetadata := account.Metadata
	newMetadata["wallets/custom_data"] = ""

	switch v := customMetadataRaw.(type) {
	case string:
		decodedMetadata := metadata.Metadata{}
		if err := json.Unmarshal([]byte(v), &decodedMetadata); err != nil {
			return errors.Wrap(err, "decoding original metadata")
		}
		for key, value := range decodedMetadata {
			newMetadata["wallets/custom_data_"+key] = value
		}
	case map[string]any:
		for k, v := range v {
			newMetadata["wallets/custom_data_"+k] = v
		}
	default:
		panic("should not happen")
	}
	account.Metadata = newMetadata

	return updateMetadata(ctx, config, account)
}
