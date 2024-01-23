package settings

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	"strconv"
	"strings"
)

func ValueOrDefault[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

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

func GetInt32(ctx core.Context, stack string, keys ...string) (*int32, error) {
	value, err := Get(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	intValue, err := strconv.ParseInt(*value, 10, 32)
	if err != nil {
		return nil, err
	}

	return pointer.For(int32(intValue)), nil
}

func GetUInt64(ctx core.Context, stack string, keys ...string) (*uint64, error) {
	value, err := Get(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	intValue, err := strconv.ParseUint(*value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &intValue, nil
}

func GetUInt16(ctx core.Context, stack string, keys ...string) (*uint16, error) {
	value, err := Get(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	intValue, err := strconv.ParseUint(*value, 10, 16)
	if err != nil {
		return nil, err
	}

	return pointer.For(uint16(intValue)), nil
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

func GetUInt(ctx core.Context, stack string, keys ...string) (*uint, error) {
	value, err := GetUInt64(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return pointer.For(uint(*value)), nil
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

func GetUInt16OrDefault(ctx core.Context, stack string, defaultValue uint16, keys ...string) (uint16, error) {
	value, err := GetUInt16(ctx, stack, keys...)
	if err != nil {
		return 0, err
	}

	if value == nil {
		return defaultValue, nil
	}
	return *value, nil
}

func GetInt32OrDefault(ctx core.Context, stack string, defaultValue int32, keys ...string) (int32, error) {
	value, err := GetInt32(ctx, stack, keys...)
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

func GetMap(ctx core.Context, stack string, keys ...string) (map[string]string, error) {
	value, err := GetString(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	ret := make(map[string]string)
	parts := strings.Split(*value, ",")
	for _, part := range parts {
		parts := strings.SplitN(part, "=", 2)
		ret[parts[0]] = parts[1]
	}

	return ret, nil
}

func GetMapOrEmpty(ctx core.Context, stack string, keys ...string) (map[string]string, error) {
	value, err := GetMap(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return map[string]string{}, nil
	}

	return value, nil
}

func GetByPriority(ctx core.Context, stack string, keys ...string) (*string, error) {
	allSettings := &v1beta1.SettingsList{}
	if err := ctx.GetClient().List(ctx, allSettings, client.MatchingFields{
		"stack":  stack,
		"keylen": fmt.Sprint(len(keys)),
	}); err != nil {
		return nil, errors.Wrap(err, "listings settings")
	}

	settings := allSettings.Items
	slices.SortFunc(settings, func(a, b v1beta1.Settings) int {
		aKeys := strings.Split(a.Spec.Key, ".")
		bKeys := strings.Split(b.Spec.Key, ".")

		for _, aKey := range aKeys {
			for _, bKey := range bKeys {
				if aKey == "*" {
					return -1
				}
				if bKey == "*" {
					return 1
				}
			}
		}

		return 0
	})

	for _, setting := range settings {
		if matchSetting(settings, keys...) {
			return &setting.Spec.Value, nil
		}
	}

	return nil, nil
}

func matchSetting(settings []v1beta1.Settings, keys ...string) bool {
	for i, setting := range settings {
		settingKeyParts := strings.Split(setting.Spec.Key, ".")
		if settingKeyParts[i] == "*" {
			continue
		}
		if settingKeyParts[i] != keys[i] {
			return false
		}
	}
	return true
}
