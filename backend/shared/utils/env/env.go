package core_utils_env

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidInputType    = errors.New("type must be a pointer to a struct")
	ErrInvalidTagFormat    = errors.New("invalid tag format, must be 'name,[required],[default]'")
	ErrRequiredFieldNoInfo = errors.New("required field has no value and no default")

	textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()
	durationType        = reflect.TypeFor[time.Duration]()
)

type objectInfo struct {
	Name       string
	Type       reflect.Type
	Required   bool
	HasDefault bool
	Default    any
}

func Process(obj any, prefix string) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return ErrInvalidInputType
	}
	return processStruct(v.Elem(), prefix)
}

func processStruct(v reflect.Value, prefix string) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if !field.IsExported() {
			continue
		}

		tag, hasTag := field.Tag.Lookup("env")

		underlying := field.Type
		isPtr := underlying.Kind() == reflect.Pointer
		if isPtr {
			underlying = underlying.Elem()
		}

		if underlying.Kind() == reflect.Struct && !implementsTextUnmarshaler(underlying) {
			nextPrefix := prefix
			if hasTag {
				info, err := getFieldInfo(tag, underlying)
				if err != nil {
					return fmt.Errorf("field %s: %w", field.Name, err)
				}
				nextPrefix = joinPrefix(prefix, info.Name)
			}

			target := fieldVal
			if isPtr {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(underlying))
				}
				target = fieldVal.Elem()
			}

			if err := processStruct(target, nextPrefix); err != nil {
				return err
			}
			continue
		}

		if !hasTag {
			continue
		}

		info, err := getFieldInfo(tag, field.Type)
		if err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}

		key := joinPrefix(prefix, info.Name)
		if err := setField(fieldVal, key, info); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	return nil
}

func setField(fieldVal reflect.Value, key string, info objectInfo) error {
	raw, ok := os.LookupEnv(key)
	if !ok {
		if info.Required {
			return fmt.Errorf("%w: %s", ErrRequiredFieldNoInfo, key)
		}
		if !info.HasDefault {
			return nil
		}
		fieldVal.Set(reflect.ValueOf(info.Default))
		return nil
	}

	val, err := parseValue(raw, info.Type)
	if err != nil {
		return fmt.Errorf("env var %s: %w", key, err)
	}
	fieldVal.Set(reflect.ValueOf(val))
	return nil
}

func implementsTextUnmarshaler(t reflect.Type) bool {
	return reflect.PointerTo(t).Implements(textUnmarshalerType)
}

func parseValue(raw string, t reflect.Type) (any, error) {
	if implementsTextUnmarshaler(t) {
		ptr := reflect.New(t)
		if err := ptr.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(raw)); err != nil {
			return nil, err
		}
		return ptr.Elem().Interface(), nil
	}
	return convertValue(raw, t)
}

func getFieldInfo(tag string, t reflect.Type) (objectInfo, error) {
	parts := strings.SplitN(tag, ",", 3)
	if parts[0] == "" {
		return objectInfo{}, ErrInvalidTagFormat
	}

	info := objectInfo{Name: parts[0], Type: t}

	if len(parts) >= 2 {
		switch parts[1] {
		case "":
		case "required":
			info.Required = true
		default:
			return objectInfo{}, ErrInvalidTagFormat
		}
	}

	if len(parts) == 3 {
		info.HasDefault = true
		def, err := parseValue(parts[2], t)
		if err != nil {
			return objectInfo{}, fmt.Errorf("failed to parse default %q as %v: %w", parts[2], t, err)
		}
		info.Default = def
	}

	return info, nil
}

func convertValue(raw string, t reflect.Type) (any, error) {
	if t == durationType {
		return time.ParseDuration(raw)
	}

	parseNum := func(parseFn func(string, int) (any, error)) (any, error) {
		v, err := parseFn(raw, t.Bits())
		if err != nil {
			return nil, fmt.Errorf("parsing %q as %s: %w", raw, t.Kind(), err)
		}
		return reflect.ValueOf(v).Convert(t).Interface(), nil
	}

	switch t.Kind() {
	case reflect.String:
		return raw, nil
	case reflect.Bool:
		return strconv.ParseBool(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseNum(func(s string, bits int) (any, error) {
			return strconv.ParseInt(s, 10, bits)
		})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseNum(func(s string, bits int) (any, error) {
			return strconv.ParseUint(s, 10, bits)
		})
	case reflect.Float32, reflect.Float64:
		return parseNum(func(s string, bits int) (any, error) {
			return strconv.ParseFloat(s, bits)
		})
	default:
		return nil, fmt.Errorf("unsupported type: %v", t.Kind())
	}
}

func joinPrefix(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "_" + name
}
