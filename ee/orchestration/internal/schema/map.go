package schema

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/pkg/errors"
)

type Context struct {
	Variables map[string]string
}

type fieldResolveError struct {
	name string
	err  error
}

func (err *fieldResolveError) Error() string {
	return fmt.Sprintf("resolving field '%s': %s", err.name, err.err)
}

func interpolate(ctx Context, v string) string {
	r := regexp.MustCompile(`\$\{[^\}]+\}`)
	return r.ReplaceAllStringFunc(v, func(key string) string {
		key = strings.TrimPrefix(key, "${")
		key = strings.TrimSuffix(key, "}")
		return ctx.Variables[key]
	})
}

func mapObjectField(ctx Context, raw any, spec reflect.Value, fieldTag tag) error {
	switch spec.Kind() {
	case reflect.Map:
		if raw == nil {
			return nil
		}
		vRaw := reflect.ValueOf(raw)
		if vRaw.Kind() != reflect.Map {
			return fmt.Errorf("expecting map, got type: %T", raw)
		}

		spec.Set(reflect.MakeMap(spec.Type()))
		for _, key := range vRaw.MapKeys() {
			value := vRaw.MapIndex(key)
			targetValue := reflect.New(spec.Type().Elem()).Elem()
			if err := mapObjectField(ctx, value.Interface(), targetValue, tag{}); err != nil {
				return err
			}
			spec.SetMapIndex(key, targetValue)
		}
	case reflect.Pointer:
		if raw == nil {
			return nil
		}
		if _, isBigInt := spec.Interface().(*big.Int); isBigInt {
			switch json := raw.(type) {
			case string:
				interpolated := interpolate(ctx, json)
				if interpolated == "" {
					interpolated = fieldTag.defaultValue
					if interpolated == "" {
						return nil
					}
				}
				bigIntValue, ok := big.NewInt(0).SetString(interpolated, 10)
				if !ok {
					return fmt.Errorf("unable to parse '%s' as big int", interpolated)
				}
				spec.Set(reflect.ValueOf(bigIntValue))
			case float64:
				spec.Set(reflect.ValueOf(big.NewInt(int64(json))))
			case nil:
				defaultValue := fieldTag.defaultValue
				if defaultValue == "" {
					return nil
				}
				bigIntValue, ok := big.NewInt(0).SetString(defaultValue, 10)
				if !ok {
					panic(fmt.Errorf("unable to parse '%s' as big int", defaultValue))
				}
				spec.Set(reflect.ValueOf(bigIntValue))
			default:
				return fmt.Errorf("expected big int or interpolated string but was %T", json)
			}
			return nil
		}
		spec.Set(reflect.New(spec.Type().Elem()))
		return mapObjectField(ctx, raw, spec.Elem(), fieldTag)
	case reflect.String:
		switch json := raw.(type) {
		case string:
			interpolated := interpolate(ctx, json)
			if interpolated == "" {
				interpolated = fieldTag.defaultValue
			}
			spec.SetString(interpolated)
		case nil:
			spec.SetString(fieldTag.defaultValue)
		default:
			return fmt.Errorf("expected string but was %T", json)
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch json := raw.(type) {
		case string:
			interpolated := interpolate(ctx, json)
			if interpolated == "" {
				interpolated = fieldTag.defaultValue
				if interpolated == "" {
					return nil
				}
			}
			uint64Value, err := strconv.ParseUint(interpolated, 10, 64)
			if err != nil {
				return fmt.Errorf("unable to resolve field '%s' to uint value", spec.Type().Name())
			}
			spec.SetUint(uint64Value)
		case float64:
			spec.SetUint(uint64(json))
		case nil:
			defaultValue := fieldTag.defaultValue
			if defaultValue == "" {
				return nil
			}
			uint64Value, err := strconv.ParseUint(defaultValue, 10, 64)
			if err != nil {
				return fmt.Errorf("unable to resolve field '%s' to uint value", spec.Type().Name())
			}
			spec.SetUint(uint64Value)
		default:
			return fmt.Errorf("expected uint or interpolated string but was %T", json)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if _, isDuration := spec.Interface().(Duration); isDuration {
			switch json := raw.(type) {
			case string:
				interpolated := interpolate(ctx, json)
				if interpolated == "" {
					interpolated = fieldTag.defaultValue
					if interpolated == "" {
						return nil
					}
				}
				duration, err := time.ParseDuration(interpolated)
				if err != nil {
					return fmt.Errorf("unable to resolve field '%s' to duration value", spec.Type().Name())
				}
				spec.SetInt(int64(duration))
			case nil:
				defaultValue := fieldTag.defaultValue
				if defaultValue == "" {
					return nil
				}
				duration, err := time.ParseDuration(defaultValue)
				if err != nil {
					return fmt.Errorf("unable to resolve field '%s' to duration value", spec.Type().Name())
				}
				spec.SetInt(int64(duration))
			default:
				return fmt.Errorf("expected uint or interpolated string but was %T", json)
			}
			return nil
		}
		switch json := raw.(type) {
		case string:
			interpolated := interpolate(ctx, json)
			if interpolated == "" {
				interpolated = fieldTag.defaultValue
				if interpolated == "" {
					return nil
				}
			}
			int64Value, err := strconv.ParseInt(interpolated, 10, 64)
			if err != nil {
				return fmt.Errorf("unable to resolve field '%s' to int value", spec.Type().Name())
			}
			spec.SetInt(int64Value)
		case float64:
			spec.SetInt(int64(json))
		case nil:
			defaultValue := fieldTag.defaultValue
			if defaultValue == "" {
				return nil
			}
			int64Value, err := strconv.ParseInt(defaultValue, 10, 64)
			if err != nil {
				return fmt.Errorf("unable to resolve field '%s' to int value", spec.Type().Name())
			}
			spec.SetInt(int64Value)
		default:
			return fmt.Errorf("expected uint or interpolated string but was %T", json)
		}
	case reflect.Bool:
		switch json := raw.(type) {
		case string:
			interpolated := strings.ToLower(interpolate(ctx, json))
			if interpolated != "true" && interpolated != "false" {
				return fmt.Errorf("unable to resolve field '%s' to bool value", spec.Type().Name())
			}
			spec.SetBool(interpolated == "true")
		case bool:
			spec.SetBool(json)
		case nil:
			defaultValue := strings.ToLower(fieldTag.defaultValue)
			if defaultValue != "true" && defaultValue != "false" {
				return fmt.Errorf("unable to resolve field '%s' to bool value", spec.Type().Name())
			}
			spec.SetBool(defaultValue == "true")
		default:
			return fmt.Errorf("expected uint or interpolated string but was %T", json)
		}
	case reflect.Float64, reflect.Float32:
		switch json := raw.(type) {
		case string:
			interpolated := strings.ToLower(interpolate(ctx, json))
			float64Value, err := strconv.ParseFloat(interpolated, 64)
			if err != nil {
				return fmt.Errorf("expected float64 or interpolated string but was %T", json)
			}
			spec.SetFloat(float64Value)
		case float64:
			spec.SetFloat(json)
		case nil:
			defaultValue := fieldTag.defaultValue
			if defaultValue == "" {
				return nil
			}
			value, err := strconv.ParseFloat(defaultValue, 64)
			if err != nil {
				return fmt.Errorf("unable to resolve field '%s' to uint value", spec.Type().Name())
			}
			spec.SetFloat(value)
		default:
			return fmt.Errorf("expected uint or interpolated string but was %T", json)
		}
	case reflect.Struct:
		if _, isDate := spec.Interface().(time.Time); isDate {
			switch json := raw.(type) {
			case string:
				interpolated := interpolate(ctx, json)
				if interpolated == "" {
					interpolated = fieldTag.defaultValue
				}
				date, err := time.ParseTime(interpolated)
				if err != nil {
					return fmt.Errorf("expected date as rfc3339 format")
				}
				spec.Set(reflect.ValueOf(date))
			case nil:
				date, err := time.ParseTime(fieldTag.defaultValue)
				if err != nil {
					return fmt.Errorf("expected date as rfc3339 format")
				}
				spec.Set(reflect.ValueOf(date))
			default:
				return fmt.Errorf("expected string but was %T", json)
			}
			return nil
		}
		asMap, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map but was %T", raw)
		}
		if err := mapObject(ctx, asMap, spec); err != nil {
			return err
		}
	}
	return nil
}

func mapObject(ctx Context, raw map[string]any, spec reflect.Value) error {
	specType := spec.Type()
	specFieldsCount := specType.NumField()
	for i := 0; i < specFieldsCount; i++ {
		specFieldType := specType.Field(i)
		specField := spec.Field(i)

		if specFieldType.Anonymous {
			if err := mapObject(ctx, raw, specField); err != nil {
				return err
			}
			continue
		}

		jsonKey := strings.Split(specFieldType.Tag.Get("json"), ",")[0]

		if err := mapObjectField(ctx, raw[jsonKey], specField, parseTag(specFieldType.Tag.Get("spec"))); err != nil {
			return &fieldResolveError{
				name: specFieldType.Name,
				err:  err,
			}
		}
	}
	return nil
}

func Resolve(ctx Context, data map[string]any, name string) (stages.Stage, error) {
	schema := stages.Get(name)
	if schema == nil {
		return nil, errors.New("unable to resolve schema " + name)
	}
	schema = reflect.New(reflect.TypeOf(schema)).Interface().(stages.Stage)
	return schema, mapObject(ctx, data, reflect.ValueOf(schema).Elem())
}
