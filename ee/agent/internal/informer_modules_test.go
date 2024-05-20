package internal_test

import (
	"encoding/json"
	"testing"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRestrictStatus(t *testing.T) {

	type testCase struct {
		incomingStatus map[string]interface{}
		expectedStatus map[string]interface{}
		expectError    bool
	}
	conditions := func() []interface{} {
		conditions := []v1beta1.Condition{}

		var count int64 = 0
		newCondition := func() v1beta1.Condition {
			count++
			return v1beta1.Condition{
				Type:               uuid.NewString(),
				Reason:             uuid.NewString(),
				Message:            uuid.NewString(),
				Status:             v1.ConditionStatus(uuid.NewString()),
				ObservedGeneration: count,
				LastTransitionTime: v1.Time{},
			}
		}

		conditions = append(conditions, newCondition())
		conditions = append(conditions, newCondition())
		return collectionutils.Map(conditions, func(c v1beta1.Condition) interface{} {
			b, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}
			var m map[string]interface{}
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatal(err)
			}
			return m
		})
	}()
	testCases := []testCase{
		{
			incomingStatus: map[string]interface{}{},
			expectedStatus: map[string]interface{}{},
			expectError:    true,
		},
		{
			incomingStatus: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectedStatus: map[string]interface{}{
				"ready": false,
			},
		},
		{
			incomingStatus: map[string]interface{}{
				"info":  "some info",
				"ready": true,
			},
			expectedStatus: map[string]interface{}{
				"info":  "some info",
				"ready": true,
			},
		},
		{
			incomingStatus: map[string]interface{}{
				"conditions": conditions,
			},

			expectedStatus: map[string]interface{}{
				"ready":      false,
				"conditions": conditions,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run("test", func(t *testing.T) {
			t.Parallel()

			status, err := internal.Restrict[v1beta1.Status](tc.incomingStatus)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, status)
		})
	}
}
