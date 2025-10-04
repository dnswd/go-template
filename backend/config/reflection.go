package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func populateWithReflection(config *Config) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue // Skip fields without env tag
		}
		
		envValue := os.Getenv(envTag)

		// Set the field value based on its type
		if err := setField(fieldValue, envValue, field.Type); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

func setField(field reflect.Value, value string, fieldType reflect.Type) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int value: %w", err)
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			field.SetUint(0)
			return nil
		}
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint value: %w", err)
		}
		field.SetUint(uintVal)
	case reflect.Bool:
		if value == "" {
			field.SetBool(false)
			return nil
		}
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool value: %w", err)
		}
		field.SetBool(boolVal)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			field.SetFloat(0)
			return nil
		}
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %w", err)
		}
		field.SetFloat(floatVal)
	default:
		// Handle custom types like Env
		if fieldType.Name() == "Env" {
			env, err := ParseEnv(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(env))
		} else {
			return fmt.Errorf("unsupported field type: %s", field.Kind())
		}
	}

	return nil
}
