package settings

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(ctx core.Context, stack string, keys ...string) (*string, error) {
	keys = Flatten(Map(keys, func(from string) []string {
		return strings.Split(from, ".")
	}))
	allSettingsTargetingStack := &v1beta1.SettingsList{}
	if err := ctx.GetClient().List(ctx, allSettingsTargetingStack, client.MatchingFields{
		"stack":  stack,
		"keylen": fmt.Sprint(len(keys)),
	}); err != nil {
		return nil, errors.Wrap(err, "listings settings")
	}

	allSettingsTargetingAllStacks := &v1beta1.SettingsList{}
	if err := ctx.GetClient().List(ctx, allSettingsTargetingAllStacks, client.MatchingFields{
		"stack":  "*",
		"keylen": fmt.Sprint(len(keys)),
	}); err != nil {
		return nil, errors.Wrap(err, "listings settings")
	}

	return findMatchingSettings(append(allSettingsTargetingStack.Items, allSettingsTargetingAllStacks.Items...), keys...)
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
		return "", core.NewMissingSettingsError(fmt.Sprintf("settings '%s' not found for stack '%s'", strings.Join(keys, "."), stack))
	}
	return *value, nil
}

func GetURL(ctx core.Context, stack string, keys ...string) (*v1beta1.URI, error) {
	value, err := GetString(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return v1beta1.ParseURL(*value)
}

func RequireURL(ctx core.Context, stack string, keys ...string) (*v1beta1.URI, error) {
	value, err := RequireString(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	return v1beta1.ParseURL(value)
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

func GetBoolOrTrue(ctx core.Context, stack string, keys ...string) (bool, error) {
	return GetBoolOrDefault(ctx, stack, true, keys...)
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

func GetAs[T any](ctx core.Context, stack string, keys ...string) (*T, error) {
	m, err := GetMap(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}

	var ret T
	ret = reflect.New(reflect.TypeOf(ret)).Elem().Interface().(T)
	if m == nil {
		return &ret, nil
	}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
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

func findMatchingSettings(settings []v1beta1.Settings, keys ...string) (*string, error) {

	// Keys can be passed as "a.b.c", instead of "a", "b", "c"
	keys = Flatten(Map(keys, func(from string) []string {
		return strings.Split(from, ".")
	}))

	slices.SortFunc(settings, sortSettingsByPriority)

	for _, setting := range settings {
		if matchSetting(setting, keys...) {
			return &setting.Spec.Value, nil
		}
	}

	return nil, nil
}

func matchSetting(setting v1beta1.Settings, keys ...string) bool {
	settingKeyParts := strings.Split(setting.Spec.Key, ".")
	for i, settingKeyPart := range settingKeyParts {
		if settingKeyPart == "*" {
			continue
		}
		if settingKeyPart != keys[i] {
			return false
		}
	}
	return true
}

func sortSettingsByPriority(a, b v1beta1.Settings) int {
	switch {
	case a.IsWildcard() && !b.IsWildcard():
		return 1
	case !a.IsWildcard() && b.IsWildcard():
		return -1
	}
	aKeys := strings.Split(a.Spec.Key, ".")
	bKeys := strings.Split(b.Spec.Key, ".")

	for i := 0; i < len(aKeys); i++ {
		if aKeys[i] == bKeys[i] {
			continue
		}
		if aKeys[i] == "*" {
			return 1
		}
		if bKeys[i] == "*" {
			return -1
		}
	}

	return 0
}

func IsTrue(v string) bool {
	return strings.ToLower(v) == "true" || v == "1"
}
