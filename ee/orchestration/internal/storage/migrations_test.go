package storage

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/stretchr/testify/require"
)

func TestMigrationsAddTemporalRunIDAsCompoundPrimaryKeyOnStages(t *testing.T) {
	db := pgtesting.NewPostgresDatabase(t)

	bunDB, err := bunconnect.OpenSQLDB(logging.TestingContext(), db.ConnectionOptions())
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
