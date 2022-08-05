package mongo

import (
	"context"
	"testing"

	"github.com/numary/webhooks-cloud/internal/env"
	"github.com/numary/webhooks-cloud/pkg/model"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_FindLastConfig(t *testing.T) {
	flagSet := pflag.NewFlagSet("TestStore", pflag.ContinueOnError)
	require.NoError(t, env.Init(flagSet))

	ctx := context.Background()
	s, err := NewConfigStore()
	require.NoError(t, err)

	require.NoError(t, s.DropConfigsCollection(ctx))

	_, err = s.InsertOneConfig(ctx, model.Config{
		Active:     true,
		EventTypes: []string{"TYPE1"},
		Endpoints:  []string{"https://www.site1.com"},
	})
	require.NoError(t, err)

	id2, err := s.InsertOneConfig(ctx, model.Config{
		Active:     true,
		EventTypes: []string{"TYPE2"},
		Endpoints:  []string{"https://www.site2.com"},
	})
	require.NoError(t, err)

	res, err := s.FindLastConfig(context.Background())
	require.NoError(t, err)

	assert.Equal(t, id2, res.ID)
}
