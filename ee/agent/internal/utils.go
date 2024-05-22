package internal

import (
	"encoding/json"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Restrict[T any](obj map[string]interface{}) (map[string]interface{}, error) {
	if len(obj) == 0 {
		return nil, errors.New("obj is empty")
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	res := new(T)
	err = json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal json in"+reflect.TypeOf(res).String())
	}

	filtered, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(filtered, &tmp); err != nil {
		return nil, err
	}

	return tmp, nil
}

func getStatus(unstructuredModule *unstructured.Unstructured) (*structpb.Struct, error) {
	status, found, err := unstructured.NestedMap(unstructuredModule.Object, "status")
	if err != nil {
		return nil, errors.Wrap(err, "unable to get status from unstructured")
	}

	if !found {
		return nil, nil
	}

	status, err = Restrict[v1beta1.Status](status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to restrict status according to v1beta1.StatusWithConditions")
	}

	protoStatus, err := structpb.NewStruct(status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert status to proto struct")
	}
	return protoStatus, nil

}
