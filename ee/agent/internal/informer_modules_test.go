package internal_test

import (
	"encoding/json"
	"testing"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestRestrictModuleStatus(t *testing.T) {

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

func TestModuleAddFunc(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	membershipClientMock := internal.NewMockMembershipClient(ctrl)
	resourceInformer := internal.NewModuleEventHandler(logging.Testing(), membershipClientMock)

	module := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": uuid.NewString(),
			},
			"status": map[string]interface{}{
				"ready": false,
			},
		},
	}

	membershipClientMock.EXPECT().Send(gomock.Any())
	resourceInformer.AddFunc(module)
	require.True(t, ctrl.Satisfied())

}

func TestModuleDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	membershipClientMock := internal.NewMockMembershipClient(ctrl)
	resourceInformer := internal.NewModuleEventHandler(logging.Testing(), membershipClientMock)

	module := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": uuid.NewString(),
			},
			"status": map[string]interface{}{
				"ready": false,
			},
		},
	}

	membershipClientMock.EXPECT().Send(gomock.Any())
	resourceInformer.DeleteFunc(module)
	require.True(t, ctrl.Satisfied())
}

func TestModuleUpdateStatusNil(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	membershipClientMock := internal.NewMockMembershipClient(ctrl)
	resourceInformer := internal.NewModuleEventHandler(logging.Testing(), membershipClientMock)

	oldModule := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": uuid.NewString(),
			},
		},
	}

	newModule := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": uuid.NewString(),
			},
		},
	}

	resourceInformer.UpdateFunc(oldModule, newModule)
	require.True(t, ctrl.Satisfied())

}
func TestModuleUpdateStatusChanged(t *testing.T) {

	type testCase struct {
		isReady  bool
		wasReady bool
	}

	testCases := []testCase{}
	for _, a := range []bool{true, false} {
		for _, b := range []bool{true, false} {
			testCases = append(testCases, testCase{a, b})
		}
	}

	for _, tc := range testCases {
		tc := tc
		t.Run("test", func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			membershipClientMock := internal.NewMockMembershipClient(ctrl)
			resourceInformer := internal.NewModuleEventHandler(logging.Testing(), membershipClientMock)

			oldModule := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": uuid.NewString(),
					},
					"status": map[string]interface{}{
						"ready": tc.wasReady,
					},
				},
			}

			newModule := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": uuid.NewString(),
					},
					"status": map[string]interface{}{
						"ready": tc.isReady,
					},
				},
			}
			if tc.isReady != tc.wasReady {
				membershipClientMock.EXPECT().Send(gomock.Any())
			}
			resourceInformer.UpdateFunc(oldModule, newModule)
			require.True(t, ctrl.Satisfied())
		})
	}

}
