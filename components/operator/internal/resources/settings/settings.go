package settings

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

func Get(ctx core.Context, stack string, keys ...string) (*string, error) {
	key := strings.Join(keys, ".")
	list := &v1beta1.SettingsList{}
	if err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		"stack": stack,
		"key":   key,
	}); err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		return nil, nil
	}
	if len(list.Items) > 1 {
		return nil, fmt.Errorf("found multiple matching setting with key '%s' and stack '%s'", key, stack)
	}

	return &list.Items[0].Spec.Value, nil
}

func GetString(ctx core.Context, stack string, keys ...string) (*string, error) {
	return Get(ctx, stack, keys...)
}

func GetStringOrDefault(ctx core.Context, stack, defaultValue string, keys ...string) (string, error) {
	value, err := GetString(ctx, stack, keys...)
	if err != nil {
		return "", err
	}
	if value == nil {
		return defaultValue, nil
	}
	return *value, nil
}

func GetStringOrEmpty(ctx core.Context, stack string, keys ...string) (string, error) {
	return GetStringOrDefault(ctx, stack, "", keys...)
}

func GetStringSlice(ctx core.Context, stack string, keys ...string) ([]string, error) {
	value, err := GetString(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, err
	}
	return strings.Split(*value, ","), nil
}

func RequireString(ctx core.Context, stack string, keys ...string) (string, error) {
	value, err := GetString(ctx, stack, keys...)
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", fmt.Errorf("settings '%s' not found for stack '%s'", strings.Join(keys, "."), stack)
	}
	return *value, nil
}

func ValueOrDefault[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

func GetInt64(ctx core.Context, stack string, keys ...string) (*int64, error) {
	value, err := Get(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	intValue, err := strconv.ParseInt(*value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &intValue, nil
}

func GetInt(ctx core.Context, stack string, keys ...string) (*int, error) {
	value, err := GetInt64(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return pointer.For(int(*value)), nil
}

func GetIntOrDefault(ctx core.Context, stack string, defaultValue int, keys ...string) (int, error) {
	value, err := GetInt(ctx, stack, keys...)
	if err != nil {
		return 0, err
	}

	if value == nil {
		return defaultValue, nil
	}
	return *value, nil
}

func GetBool(ctx core.Context, stack string, keys ...string) (*bool, error) {
	value, err := Get(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return pointer.For(*value == "true"), nil
}

func GetBoolOrDefault(ctx core.Context, stack string, defaultValue bool, keys ...string) (bool, error) {
	value, err := GetBool(ctx, stack, keys...)
	if err != nil {
		return false, err
	}
	if value == nil {
		return defaultValue, nil
	}
	return *value, nil
}

func GetBoolOrFalse(ctx core.Context, stack string, keys ...string) (bool, error) {
	return GetBoolOrDefault(ctx, stack, false, keys...)
}
