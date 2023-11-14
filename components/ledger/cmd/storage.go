package cmd

import (
	"context"
	"errors"

	"github.com/formancehq/ledger/internal/storage"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewBucket() *cobra.Command {
	return &cobra.Command{
		Use: "buckets",
	}
}

func NewBucketInit() *cobra.Command {
	cmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := service.New(
				cmd.OutOrStdout(),
				resolveOptions(
					cmd.OutOrStdout(),
					fx.Invoke(func(storageDriver *driver.Driver, lc fx.Lifecycle) {
						lc.Append(fx.Hook{
							OnStart: func(ctx context.Context) error {
								name := viper.GetString("name")
								if name == "" {
									return errors.New("name is empty")
								}

								exists, err := storageDriver.GetSystemStore().ExistsLedger(ctx, name)
								if err != nil {
									return err
								}

								if exists {
									return errors.New("ledger already exists")
								}

								store, err := storageDriver.GetLedgerStore(ctx, name)
								if err != nil {
									return err
								}

								return store.Initialize(ctx)
							},
						})
					}))...,
			)
			return app.Start(cmd.Context())
		},
	}
	cmd.Flags().String("name", "default", "Ledger name")
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}
	return cmd
}

func NewBucketList() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := service.New(cmd.OutOrStdout(),
				resolveOptions(
					cmd.OutOrStdout(),
					fx.Invoke(func(storageDriver *driver.Driver, lc fx.Lifecycle) {
						lc.Append(fx.Hook{
							OnStart: func(ctx context.Context) error {
								ledgers, err := storageDriver.GetSystemStore().ListLedgers(ctx)
								if err != nil {
									return err
								}
								if len(ledgers) == 0 {
									logging.FromContext(ctx).Info("No ledger found.")
									return nil
								}
								logging.FromContext(ctx).Infof("Ledgers: %v", ledgers)
								return nil
							},
						})
					}),
				)...,
			)
			return app.Start(cmd.Context())
		},
	}
	return cmd
}

func NewBucketUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "upgrade",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			driver := driver.New(storage.ConnectionOptionsFromFlags(viper.GetViper(), cmd.OutOrStdout(), viper.GetBool(service.DebugFlag)))
			if err := driver.Initialize(cmd.Context()); err != nil {
				return err
			}
			defer func() {
				_ = driver.Close()
			}()

			name := args[0]

			bucket, err := driver.GetBucket(cmd.Context(), name)
			if err != nil {
				return err
			}

			logger := service.GetDefaultLogger(cmd.OutOrStdout(), viper.GetBool(service.DebugFlag), false)

			return bucket.Migrate(logging.ContextWithLogger(cmd.Context(), logger))
		},
	}
	return cmd
}

func NewBucketUpgradeAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "upgrade-all",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			logger := service.GetDefaultLogger(cmd.OutOrStdout(), viper.GetBool(service.DebugFlag), false)
			ctx := logging.ContextWithLogger(cmd.Context(), logger)

			driver := driver.New(storage.ConnectionOptionsFromFlags(viper.GetViper(), cmd.OutOrStdout(), viper.GetBool(service.DebugFlag)))
			defer func() {
				_ = driver.Close()
			}()

			if err := driver.Initialize(ctx); err != nil {
				return err
			}

			return driver.UpgradeAllBuckets(ctx)
		},
	}
	return cmd
}

func NewBucketDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := service.New(
				cmd.OutOrStdout(),
				resolveOptions(
					cmd.OutOrStdout(),
					fx.Invoke(func(storageDriver *driver.Driver, lc fx.Lifecycle) {
						lc.Append(fx.Hook{
							OnStart: func(ctx context.Context) error {
								name := args[0]
								store, err := storageDriver.GetLedgerStore(ctx, name)
								if err != nil {
									return err
								}

								if err := store.Delete(ctx); err != nil {
									return err
								}
								return nil
							},
						})
					}))...,
			)
			return app.Start(cmd.Context())
		},
	}
	return cmd
}
