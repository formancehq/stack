package settings

import (
	"github.com/formancehq/go-libs/collectionutils"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func GetResourceRequirements(ctx core.Context, stack string, keys ...string) (*v1.ResourceRequirements, error) {
	limits, err := GetResourceList(ctx, stack, append(keys, "limits")...)
	if err != nil {
		return nil, err
	}

	requests, err := GetResourceList(ctx, stack, append(keys, "requests")...)
	if err != nil {
		return nil, err
	}

	claims, err := GetStringSlice(ctx, stack, append(keys, "claims")...)
	if err != nil {
		return nil, err
	}

	return &v1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
		Claims: collectionutils.Map(claims, func(from string) v1.ResourceClaim {
			return v1.ResourceClaim{
				Name: from,
			}
		}),
	}, nil
}

func GetResourceList(ctx core.Context, stack string, keys ...string) (v1.ResourceList, error) {
	value, err := GetMap(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	ret := v1.ResourceList{}
	for key, qty := range value {
		ret[v1.ResourceName(key)], err = resource.ParseQuantity(qty)
	}

	return ret, nil
}
