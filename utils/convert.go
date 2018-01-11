package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}

func String2Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Map2Struct(m map[string]interface{}, v interface{}) error {
	for k, v := range m {
		structValue := reflect.ValueOf(v).Elem()
		structFieldValue := structValue.FieldByName(k)

		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", k)
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", k)
		}

		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(v)
		if structFieldType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)

	}

	return nil
}
