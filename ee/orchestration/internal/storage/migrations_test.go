package storage

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/uptrace/bun"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/require"
)

func TestMigrationsAddTemporalRunIDAsCompoundPrimaryKeyOnStages(t *testing.T) {
	db := srv.NewDatabase()

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	bunDB, err := bunconnect.OpenSQLDB(logging.TestingContext(), db.ConnectionOptions(), hooks...)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, bunDB.Close())
	}()
	require.NoError(t, MigrateUntil(logging.TestingContext(), bunDB, 5))

	_, err = bunDB.Exec(`
		insert into workflows(config, id)
		values ('{}'::jsonb, '1');
		
		insert into workflow_instances(workflow_id, id)
		values ('1', '1');

		insert into workflow_instance_stage_statuses(instance_id, stage)
		values('1', 0);
	`)
	require.NoError(t, err)

	// The migration 6 will update primary key of the table `workflow_instance_stage_statuses`
	// Since we insert some data before, this test ensure the schema is not breaking after the migration
	require.NoError(t, Migrate(logging.TestingContext(), bunDB))
}
